package orm_framework

import (
	"context"
	"github.com/borntodie-new/orm-framework/internal/errs"
	"strings"
)

// UpdateSQL 修改语句的原型
// 1. 需要实现 Executer 接口，用于执行修改SQL语句的功能
// 2. 需要实现 Builder 接口，用于构建SQL语句和保存SQL的参数
type UpdateSQL[T any] struct {
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
	// values 需要修改的数据
	values map[string]any
}

func (u *UpdateSQL[T]) Where(condition ...Predicate) *UpdateSQL[T] {
	u.where = append(u.where, condition...)
	return u
}

func (u *UpdateSQL[T]) Table(tableName string) *UpdateSQL[T] {
	u.table = tableName
	return u
}

func (u *UpdateSQL[T]) Values(fieldName string, data any) *UpdateSQL[T] {
	if u.values == nil {
		u.values = make(map[string]any)
	}
	u.values[fieldName] = data
	return u
}

// addArgs 添加SQL参数
func (u *UpdateSQL[T]) addArgs(val any) {
	if val == nil {
		return
	}
	u.args = append(u.args, val)
}

// buildWhere 构建 WHERE 语句
func (u *UpdateSQL[T]) buildWhere() error {
	if len(u.where) <= 0 {
		return nil
	}
	u.sb.WriteString(" WHERE ")
	p := u.where[0]
	for i := 1; i < len(u.where)-1; i++ {
		p.AND(u.where[i])
	}
	return u.buildFields(p)
}

// buildFields 构建WHERE语句
func (u *UpdateSQL[T]) buildFields(exp Expression) error {
	switch typ := exp.(type) {
	case nil:
		return nil
	case Field:
		var t T
		m, err := u.db.manager.Get(t)
		if err != nil {
			return err
		}
		// 这是纯字段
		// 注意 Field传入的是Go中的字段名，设置到SQL上的是SQL中的列名
		u.sb.WriteByte('(')
		u.sb.WriteByte('`')
		fd, ok := m.FieldsMap[typ.fieldName]
		if !ok {
			return errs.NewErrNotSupportUnknownField(typ.fieldName)
		}
		u.sb.WriteString(fd.ColumnName)
		u.sb.WriteByte('`')
	case Predicate:
		// 这里需要递归实现，因为是 Predicate 类型，可能是 Field 也可能是 Value

		// 构建左边
		if err := u.buildFields(typ.left); err != nil {
			return err
		}
		// 构建操作类型
		u.sb.WriteString(typ.op.String())
		// 构建右边
		if err := u.buildFields(typ.right); err != nil {
			return err
		}
	case Value:
		// 这里是字段值
		u.sb.WriteString("?")
		u.addArgs(typ.val)
		u.sb.WriteByte(')')
	default:
		return errs.ErrNotSupportPredicate
	}
	return nil
}

// buildValues 构建 赋值 子句
func (u *UpdateSQL[T]) buildValues() error {
	if len(u.values) <= 0 {
		return errs.ErrNotUpdateSQLSetClause
	}
	idx := 0
	var t T
	m, err := u.db.manager.Get(t)
	if err != nil {
		return err
	}
	for fieldName, value := range u.values {
		if idx > 0 {
			u.sb.WriteString(", ")
		}
		fd, ok := m.FieldsMap[fieldName]
		if !ok {
			return errs.NewErrNotSupportUnknownField(fieldName)
		}
		// 设置列名
		u.sb.WriteByte('`')
		u.sb.WriteString(fd.ColumnName)
		u.sb.WriteByte('`')
		// 设置占位符
		u.sb.WriteString(" = ?")
		// 保存数据
		u.addArgs(value)
		idx++
	}
	return nil
}

// ExecuteWithContext 执行SQL语句
func (u *UpdateSQL[T]) ExecuteWithContext(ctx context.Context) (*Result, error) {
	sqlInfo, err := u.Build()
	if err != nil {
		return nil, err
	}
	res, err := u.db.db.ExecContext(ctx, sqlInfo.SQL, sqlInfo.Args...)
	if err != nil {
		return &Result{
			err: err,
			res: nil,
		}, nil
	}
	return &Result{res: res}, err
}

// Build 构造SQL语句和维护SQL参数
// UPDATE `test_model` SET `first_name` = 'Fred' WHERE `id` = 1;
func (u *UpdateSQL[T]) Build() (*SQLInfo, error) {
	// 构建SQL基本架构
	u.sb.WriteString("UPDATE ")
	var t T
	m, err := u.db.manager.Get(t)
	if err != nil {
		return nil, err
	}
	// 构建表名
	u.sb.WriteByte('`')
	if u.table != "" {
		u.sb.WriteString(u.table)
	} else {
		u.sb.WriteString(m.TableName)
	}
	u.sb.WriteByte('`')
	u.sb.WriteString(" SET ")
	// TODO 构建赋值子句
	if err = u.buildValues(); err != nil {
		return nil, err
	}
	// 构建WHERE语句
	if err = u.buildWhere(); err != nil {
		return nil, err
	}
	u.sb.WriteByte(';')
	res := &SQLInfo{SQL: u.sb.String(), Args: u.args}
	return res, nil
}

func NewUpdateSQL[T any](db *DB) *UpdateSQL[T] {
	return &UpdateSQL[T]{
		sb:   &strings.Builder{},
		args: []any{},
		db:   db,
	}
}
