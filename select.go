package orm_framework

import (
	"context"
	"database/sql"
	"github.com/borntodie-new/orm-framework/internal/errs"
	"github.com/borntodie-new/orm-framework/model"
	"reflect"
)

// SelectSQL 查询语句
// 1. 需要实现 Builder 接口，用于构建SQL语句和保存SQL的参数
// 2. 需要实现 Querier 接口，用于接收SQL返回的结果集
type SelectSQL[T any] struct {
	// sb 构建SQL语句的，性能好
	// sb *strings.Builder
	// where SQL 中的 WHERE 语句
	where []Predicate
	// args SQL语句中的参数
	// args []any
	// db 全局的、自定义的数据库连接对象
	db *DB
	// fields 查询字段
	fields []string

	// model 在语句层面维护表模型
	// model *model.Model
	// builder 抽象出新的 SQL 构造器
	*builder
}

func (s *SelectSQL[T]) Where(condition ...Predicate) *SelectSQL[T] {
	s.where = append(s.where, condition...)
	return s
}

func (s *SelectSQL[T]) Fields(fields ...string) *SelectSQL[T] {
	s.fields = append(s.fields, fields...)
	return s
}

func (s *SelectSQL[T]) addArgs(val any) {
	if val == nil {
		return
	}
	s.args = append(s.args, val)
}

// buildWhere 构建 WHERE 语句
func (s *SelectSQL[T]) buildWhere() error {
	if len(s.where) <= 0 {
		return nil
	}
	s.sb.WriteString(" WHERE ")
	p := s.where[0]
	for i := 1; i <= len(s.where)-1; i++ {
		p = p.AND(s.where[i])
	}
	return s.buildFields(p)
}

// buildFields 构建WHERE语句
func (s *SelectSQL[T]) buildFields(exp Expression) error {
	switch typ := exp.(type) {
	case nil:
		return nil
	case Field:
		//var t T
		//m, err := s.db.manager.Get(t)
		//if err != nil {
		//	return err
		//}
		// 这是纯字段
		// 注意 Field传入的是Go中的字段名，设置到SQL上的是SQL中的列名
		s.sb.WriteByte('(')
		s.sb.WriteByte('`')
		fd, ok := s.model.FieldsMap[typ.fieldName]
		if !ok {
			return errs.NewErrNotSupportUnknownField(typ.fieldName)
		}
		s.sb.WriteString(fd.ColumnName)
		s.sb.WriteByte('`')
	case Predicate:
		// 这里需要递归实现，因为是 Predicate 类型，可能是 Field 也可能是 Value

		// 构建左边
		if err := s.buildFields(typ.left); err != nil {
			return err
		}
		// 构建操作类型
		s.sb.WriteString(typ.op.String())
		// 构建右边
		if err := s.buildFields(typ.right); err != nil {
			return err
		}
	case Value:
		// 这里是字段值
		s.sb.WriteString("?")
		s.addArgs(typ.val)
		s.sb.WriteByte(')')
	default:
		return errs.ErrNotSupportPredicate
	}
	return nil
}

func (s *SelectSQL[T]) setFields(res *sql.Rows) (*T, error) {
	// 最终的结果
	tp := new(T)
	//var t T
	//m, err := s.db.manager.Get(t)
	//if err != nil {
	//	return nil, err
	//}
	// 重头戏——如何将SQL的结果集映射成Go中的struct
	orderColumnsStr, err := res.Columns()
	if err != nil {
		return nil, err
	}
	// val := reflect.Indirect(reflect.ValueOf(new(T)))
	val := reflect.ValueOf(tp).Elem()
	receiptFields := make([]any, 0, len(s.model.Fields))                    // 保存Scan需要数据
	receiptInterfaceFields := make([]reflect.Value, 0, len(s.model.Fields)) // 用于保存每个字段的 Value类型
	for _, str := range orderColumnsStr {
		fd, ok := s.model.ColumnsMap[str]
		if !ok {
			return nil, errs.NewErrNotSupportUnknownColumn(str)
		}
		temp := reflect.New(fd.Type)
		receiptFields = append(receiptFields, temp.Interface())
		receiptInterfaceFields = append(receiptInterfaceFields, temp.Elem())
	}
	// 接收SQL返回的结果数据
	err = res.Scan(receiptFields...)
	if err != nil {
		return nil, err
	}
	// 将Scan出来的数据设置到 tp 结构体字段上
	for idx, str := range orderColumnsStr {
		fd, ok := s.model.ColumnsMap[str]
		if !ok {
			return nil, errs.NewErrNotSupportUnknownColumn(str)
		}
		val.FieldByName(fd.FieldName).Set(receiptInterfaceFields[idx])
	}
	return tp, nil
}

// QueryWithContext 查询多条数据
func (s *SelectSQL[T]) QueryWithContext(ctx context.Context) ([]*T, error) {
	// 获取 SQL 语句 和 SQL 参数
	sqlInfo, err := s.Build()
	if err != nil {
		return nil, err
	}
	// 执行 SQL 语句
	res, err := s.db.db.QueryContext(ctx, sqlInfo.SQL, sqlInfo.Args...)
	if err != nil {
		return nil, err
	}
	tps := make([]*T, 0)
	for res.Next() {
		tp, err := s.setFields(res)
		if err != nil {
			return nil, err
		}
		tps = append(tps, tp)
	}
	return tps, nil
}

// QueryRawWithContext 查询单条数据
// 这里注意一下哈：这是查询单条记录的，但我们内部使用的是查询多条的API
// 但是但是，我们如果只Scan一次，就表示我们只获取第一条数据
func (s *SelectSQL[T]) QueryRawWithContext(ctx context.Context) (*T, error) {
	// 获取 SQL 语句 和 SQL 参数
	sqlInfo, err := s.Build()
	if err != nil {
		return nil, err
	}
	// 执行 SQL 语句
	res, err := s.db.db.QueryContext(ctx, sqlInfo.SQL, sqlInfo.Args...)
	if err != nil {
		return nil, err
	}
	if !res.Next() {
		return nil, errs.ErrNoRows
	}
	//if !res.Next() {
	//	return nil, errs.ErrNoRows
	//}
	//// 最终的结果
	//tp := new(T)
	////var t T
	////m, err := s.db.manager.Get(t)
	////if err != nil {
	////	return nil, err
	////}
	//// 重头戏——如何将SQL的结果集映射成Go中的struct
	//orderColumnsStr, err := res.Columns()
	//if err != nil {
	//	return nil, err
	//}
	//// val := reflect.Indirect(reflect.ValueOf(new(T)))
	//val := reflect.ValueOf(tp).Elem()
	//receiptFields := make([]any, 0, len(s.model.Fields))                    // 保存Scan需要数据
	//receiptInterfaceFields := make([]reflect.Value, 0, len(s.model.Fields)) // 用于保存每个字段的 Value类型
	//for _, str := range orderColumnsStr {
	//	fd, ok := s.model.ColumnsMap[str]
	//	if !ok {
	//		return nil, errs.NewErrNotSupportUnknownColumn(str)
	//	}
	//	temp := reflect.New(fd.Type)
	//	receiptFields = append(receiptFields, temp.Interface())
	//	receiptInterfaceFields = append(receiptInterfaceFields, temp.Elem())
	//}
	//// 接收SQL返回的结果数据
	//err = res.Scan(receiptFields...)
	//if err != nil {
	//	return nil, err
	//}
	//// 将Scan出来的数据设置到 tp 结构体字段上
	//for idx, str := range orderColumnsStr {
	//	fd, ok := s.model.ColumnsMap[str]
	//	if !ok {
	//		return nil, errs.NewErrNotSupportUnknownColumn(str)
	//	}
	//	val.FieldByName(fd.FieldName).Set(receiptInterfaceFields[idx])
	//}
	//
	//return tp, nil
	return s.setFields(res)
}

// buildColumns 构建字段
// 功能作用和 InsertSQL 中的 buildFields 功能一样，只不过在 SelectSQL 中已经有一个 buildFields 方法了
func (s *SelectSQL[T]) buildColumns() error {
	//var t T
	//m, err := s.db.manager.Get(t)
	//if err != nil {
	//	return err
	//}
	orderFields := make([]*model.Field, 0, len(s.model.Fields))
	// 处理用户自定义字段情况
	for _, fieldName := range s.fields {
		fd, ok := s.model.FieldsMap[fieldName]
		if !ok {
			return errs.NewErrNotSupportUnknownField(fieldName)
		}
		orderFields = append(orderFields, fd)
	}
	// 处理用户没有指定字段情况
	if len(s.fields) == 0 {
		orderFields = s.model.Fields
	}
	for idx, field := range orderFields {
		if idx > 0 {
			s.sb.WriteString(", ")
		}
		s.sb.WriteByte('`')
		s.sb.WriteString(field.ColumnName)
		s.sb.WriteByte('`')
	}
	return nil
}

func (s *SelectSQL[T]) Build() (*SQLInfo, error) {
	s.sb.WriteString("SELECT ")
	// 获取表模型
	var err error
	s.model, err = s.db.manager.Get(new(T))
	if err != nil {
		return nil, err
	}
	// TODO 构建查询字段
	if err = s.buildColumns(); err != nil {
		return nil, err
	}
	s.sb.WriteString(" FROM ")
	// 构建表名
	s.sb.WriteByte('`')
	s.sb.WriteString(s.model.TableName)
	s.sb.WriteByte('`')

	// 构建 WHERE 子句
	if err = s.buildWhere(); err != nil {
		return nil, err
	}
	s.sb.WriteByte(';')
	res := &SQLInfo{SQL: s.sb.String(), Args: s.args}
	return res, nil
}

// NewSelectSQL 初始化SELECT语句对象
func NewSelectSQL[T any](db *DB) *SelectSQL[T] {
	return &SelectSQL[T]{
		// sb:   &strings.Builder{},
		// args: []any{},
		builder: newBuilder(),
		db:      db,
	}
}
