package orm_framework

import (
	"context"
	"github.com/borntodie-new/orm-framework/internal/errs"
	"github.com/borntodie-new/orm-framework/internal/valuer"
	"github.com/borntodie-new/orm-framework/model"
)

// 仅支持原生的 SELECT 语句

// RawSQL 原生查询语句
// 1. 需要实现 Builder 接口，用于构建 SQL 语句。只是走一个流程，不会真的构建的。为了统一性而已
// 2. 需要实现 Querier 接口，用于执行 SELECT 语句
type RawSQL[T any] struct {
	// sql 原生 SQL 语句
	sql string
	// args SQL 所需的参数
	args []any
	// db 全局的、自定义的数据库连接对象
	db *DB
	// models 维护一个表模型
	model *model.Model
	// valuer 公共映射值方法
	valuer valuer.FactoryValuer
}

//func (r *RawSQL[T]) setFields(res *sql.Rows) (*T, error) {
//	// 最终的结果
//	tp := new(T)
//	//var t T
//	//m, err := s.db.manager.Get(t)
//	//if err != nil {
//	//	return nil, err
//	//}
//	// 重头戏——如何将SQL的结果集映射成Go中的struct
//	orderColumnsStr, err := res.Columns()
//	if err != nil {
//		return nil, err
//	}
//	// val := reflect.Indirect(reflect.ValueOf(new(T)))
//	val := reflect.ValueOf(tp).Elem()
//	receiptFields := make([]any, 0, len(r.model.Fields))                    // 保存Scan需要数据
//	receiptInterfaceFields := make([]reflect.Value, 0, len(r.model.Fields)) // 用于保存每个字段的 Value类型
//	for _, str := range orderColumnsStr {
//		fd, ok := r.model.ColumnsMap[str]
//		if !ok {
//			return nil, errs.NewErrNotSupportUnknownColumn(str)
//		}
//		temp := reflect.New(fd.Type)
//		receiptFields = append(receiptFields, temp.Interface())
//		receiptInterfaceFields = append(receiptInterfaceFields, temp.Elem())
//	}
//	// 接收SQL返回的结果数据
//	err = res.Scan(receiptFields...)
//	if err != nil {
//		return nil, err
//	}
//	// 将Scan出来的数据设置到 tp 结构体字段上
//	for idx, str := range orderColumnsStr {
//		fd, ok := r.model.ColumnsMap[str]
//		if !ok {
//			return nil, errs.NewErrNotSupportUnknownColumn(str)
//		}
//		val.FieldByName(fd.FieldName).Set(receiptInterfaceFields[idx])
//	}
//	return tp, nil
//}

func (r *RawSQL[T]) QueryWithContext(ctx context.Context) ([]*T, error) {
	// 获取 SQL 语句 和 SQL 参数
	sqlInfo, err := r.Build()
	if err != nil {
		return nil, err
	}
	// 执行 SQL 语句
	res, err := r.db.db.QueryContext(ctx, sqlInfo.SQL, sqlInfo.Args...)
	if err != nil {
		return nil, err
	}
	tps := make([]*T, 0)
	for res.Next() {
		tp := new(T)
		val := r.valuer(r.model, tp)
		err = val.SetField(res)
		if err != nil {
			return nil, err
		}
		tps = append(tps, tp)
	}
	return tps, nil
}

func (r *RawSQL[T]) QueryRawWithContext(ctx context.Context) (*T, error) {
	// 获取 SQL 语句 和 SQL 参数
	sqlInfo, err := r.Build()
	if err != nil {
		return nil, err
	}
	// 执行 SQL 语句
	res, err := r.db.db.QueryContext(ctx, sqlInfo.SQL, sqlInfo.Args...)
	if err != nil {
		return nil, err
	}
	if !res.Next() {
		return nil, errs.ErrNoRows
	}
	tp := new(T)
	val := r.valuer(r.model, tp)
	err = val.SetField(res)
	return tp, err
}

func (r *RawSQL[T]) Build() (*SQLInfo, error) {
	var err error
	r.model, err = r.db.manager.Get(new(T))
	if r.sql == "" {
		return nil, errs.ErrNoSQL
	}
	return &SQLInfo{
		SQL:  r.sql,
		Args: r.args,
	}, err
}

func NewRawSQL[T any](db *DB, valuer valuer.FactoryValuer, sql string, args ...any) *RawSQL[T] {
	// 为什么不在这里将 model 初始化好？
	// 为了不打断我们链式调用，因为获取 model 可能会出现错误，如果将 error 返回，就会打断链式调用
	return &RawSQL[T]{
		sql:    sql,
		args:   args,
		db:     db,
		valuer: valuer,
	}
}
