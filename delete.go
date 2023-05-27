package orm_framework

import (
	"context"
	"database/sql"
	"github.com/borntodie-new/orm-framework/internal/errs"
	"reflect"
	"strings"
)

// DeleteSQL 删除语句的原型
// 1. 需要实现 Executer 接口，用于执行删除SQL语句功能
// 2. 需要实现 Builder 接口，用于构建SQL语句和保存SQL的参数
type DeleteSQL[T any] struct {
	sb    *strings.Builder
	table string
	where []Predicate
	args  []any
}

// Where 设置SQL的执行条件
// 在 DELETE 语句中有点特殊，就是说，DELETE必须带上WHERE条件，否则就会把整张表的数据都删除的，切记切记
func (d *DeleteSQL[T]) Where(condition ...Predicate) *DeleteSQL[T] {
	// d.where = condition 这样操作也是可以的
	d.where = append(d.where, condition...)
	return d
}

// Table 设置模型名字
// 因为在Go对于结构体的命名规范可能和SQL中的命名规范不一样，所以我们需要显性的设置
// 如果用户没有显性指定，那我们就按照框架的默认形式为模型设置名字
func (d *DeleteSQL[T]) Table(tableName string) *DeleteSQL[T] {
	d.table = tableName
	return d
}

// Build 构建SQL语句
func (d *DeleteSQL[T]) Build() (*SQLInfo, error) {
	// 构建 DELETE 基本框架
	d.sb.WriteString("DELETE FROM ")
	// 构建 DELETE 的表名

	if d.table == "" {
		var t T
		typ := reflect.TypeOf(t)
		d.sb.WriteByte('`')
		d.sb.WriteString(typ.Name())
		d.sb.WriteByte('`')
	} else {
		d.sb.WriteByte('`')
		d.sb.WriteString(d.table)
		d.sb.WriteByte('`')
	}
	// 构建 WHERE 语句
	if err := d.buildWhere(); err != nil {
		return nil, err
	}
	d.sb.WriteByte(';')
	res := &SQLInfo{SQL: d.sb.String(), Args: d.args}
	return res, nil
}

func (d *DeleteSQL[T]) buildWhere() error {
	if len(d.where) <= 0 {
		return nil
	}
	d.sb.WriteString(" WHERE ")
	p := d.where[0]
	for i := 1; i < len(d.where)-1; i++ {
		p.AND(d.where[i])
	}
	return d.buildFields(p)
}

// buildFields 构建WHERE语句
func (d *DeleteSQL[T]) buildFields(exp Expression) error {
	switch typ := exp.(type) {
	case nil:
		return nil
	case Field:
		// 这是纯字段
		d.sb.WriteByte('(')
		d.sb.WriteByte('`')
		d.sb.WriteString(typ.fieldName)
		d.sb.WriteByte('`')
	case Predicate:
		// 这里需要递归实现，因为是 Predicate 类型，可能是 Field 也可能是 Value

		// 构建左边
		if err := d.buildFields(typ.left); err != nil {
			return err
		}
		// 构建操作类型
		d.sb.WriteString(typ.op.String())
		// 构建右边
		if err := d.buildFields(typ.right); err != nil {
			return err
		}
	case Value:
		// 这里是字段值
		d.sb.WriteString("?")
		d.addArgs(typ.val)
		d.sb.WriteByte(')')
	default:
		return errs.ErrNotSupportPredicate
	}
	return nil
}

func (d *DeleteSQL[T]) addArgs(val any) {
	if val == nil {
		return
	}
	d.args = append(d.args, val)
}

// ExecuteWithContext 执行SQL语句
func (d *DeleteSQL[T]) ExecuteWithContext(ctx context.Context) (sql.Result, error) {
	//TODO implement me
	panic("implement me")
}

// NewDeleteSQL 这是初始化一个 DeleteSQL 对象
// 并且希望能够通过链式调用来使用
func NewDeleteSQL[T any]() *DeleteSQL[T] {
	return &DeleteSQL[T]{
		sb:   &strings.Builder{},
		args: []any{},
	}
}

/*
我们预期是怎么使用这个 DeleteSQL 的
NewDeleteSQL[*TestModel]().Build()
DELETE FROM `test_model` where `id` = 1;
*/
