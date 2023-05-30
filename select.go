package orm_framework

import (
	"context"
	"github.com/borntodie-new/orm-framework/internal/errs"
	"github.com/borntodie-new/orm-framework/model"
	"strings"
)

// SelectSQL 查询语句
// 1. 需要实现 Builder 接口，用于构建SQL语句和保存SQL的参数
// 2. 需要实现 Querier 接口，用于接收SQL返回的结果集
type SelectSQL[T any] struct {
	// sb 构建SQL语句的，性能好
	sb *strings.Builder
	// table 模型名 || 结构体名字
	table string
	// where SQL 中的 WHERE 语句
	where []Predicate
	// args SQL语句中的参数
	args []any
	// db 全局的、自定义的数据库连接对象
	db *DB
	// fields 查询字段
	fields []string
}

func (s *SelectSQL[T]) Where(condition ...Predicate) *SelectSQL[T] {
	s.where = append(s.where, condition...)
	return s
}

func (s *SelectSQL[T]) Table(tableName string) *SelectSQL[T] {
	s.table = tableName
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
		var t T
		m, err := s.db.manager.Get(t)
		if err != nil {
			return err
		}
		// 这是纯字段
		// 注意 Field传入的是Go中的字段名，设置到SQL上的是SQL中的列名
		s.sb.WriteByte('(')
		s.sb.WriteByte('`')
		fd, ok := m.FieldsMap[typ.fieldName]
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

// QueryWithContext 查询多条数据
func (s *SelectSQL[T]) QueryWithContext(ctx context.Context) (*T, error) {
	//TODO implement me
	panic("implement me")
}

// QueryRawWithContext 查询单条数据
func (s *SelectSQL[T]) QueryRawWithContext(ctx context.Context) ([]*T, error) {
	//TODO implement me
	panic("implement me")
}

// buildColumns 构建字段
// 功能作用和 InsertSQL 中的 buildFields 功能一样，只不过在 SelectSQL 中已经有一个 buildFields 方法了
func (s *SelectSQL[T]) buildColumns() error {
	var t T
	m, err := s.db.manager.Get(t)
	if err != nil {
		return err
	}
	orderFields := make([]*model.Field, 0, len(m.Fields))
	// 处理用户自定义字段情况
	for _, fieldName := range s.fields {
		fd, ok := m.FieldsMap[fieldName]
		if !ok {
			return errs.NewErrNotSupportUnknownField(fieldName)
		}
		orderFields = append(orderFields, fd)
	}
	// 处理用户没有指定字段情况
	if len(s.fields) == 0 {
		orderFields = m.Fields
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
	var t T
	m, err := s.db.manager.Get(t)
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
	if s.table != "" {
		s.sb.WriteString(s.table)
	} else {
		s.sb.WriteString(m.TableName)
	}
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
		sb:   &strings.Builder{},
		args: []any{},
		db:   db,
	}
}
