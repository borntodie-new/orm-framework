package orm_framework

import (
	"context"
	"github.com/borntodie-new/orm-framework/internal/errs"
	"github.com/borntodie-new/orm-framework/model"
	"reflect"
)

// InsertSQL 修改语句的原型
// 1. 需要实现 Executer 接口，用于执行修改SQL语句的功能
// 2. 需要实现 Builder 接口，用于构建SQL语句和保存SQL的参数
type InsertSQL[T any] struct {
	// sb 构建SQL语句的，性能好
	// sb *strings.Builder
	// args SQL语句中的参数
	// args []any
	// db 全局的、自定义的数据库连接对象
	db *DB
	// values 插入的具体的值
	values []T
	// fields 指定需要插入的字段名 Go 中的
	fields []string
	// model 维护一个表模型
	// model *model.Model
	// builder 抽象出新的 SQL 构造器
	*builder
}

func (i *InsertSQL[T]) Values(values ...T) *InsertSQL[T] {
	i.values = append(i.values, values...)
	return i
}

func (i *InsertSQL[T]) Fields(fields ...string) *InsertSQL[T] {
	i.fields = append(i.fields, fields...)
	return i
}

// addArgs 添加SQL参数
func (i *InsertSQL[T]) addArgs(val any) {
	if val == nil {
		return
	}
	i.args = append(i.args, val)
}

// buildValues 构建 VALUES 子句
// 该函数的重要功能如下
// 1. 构建 len(i.values)个(?,?,?,...)
// 2. 将i.values中的所有字段添加到i.args中
func (i *InsertSQL[T]) buildColumnsAndValues() error {
	if len(i.values) <= 0 {
		return errs.ErrNotInsertSQLValuesClause
	}
	orderFields := make([]*model.Field, 0, len(i.model.Fields))
	// 将用户指定的字段信息添加到排好序的 orderFields 切片中
	for _, fieldName := range i.fields {
		field, ok := i.model.FieldsMap[fieldName]
		if !ok {
			return errs.NewErrNotSupportUnknownField(fieldName)
		}
		orderFields = append(orderFields, field)
	}
	// 如果用户没有指定字段顺序，就用默认的
	if len(i.fields) == 0 {
		orderFields = i.model.Fields
	}

	// 构建具体的列名
	i.sb.WriteByte('(')
	for idx, field := range orderFields {
		if idx > 0 {
			i.sb.WriteString(", ")
		}
		i.sb.WriteByte('`')
		i.sb.WriteString(field.ColumnName)
		i.sb.WriteByte('`')
	}
	i.sb.WriteByte(')')

	// 构建占位符
	// len(orderFields)*len(i.values) 计算出要有多少个参数，就有多少个?占位符
	i.sb.WriteString(" VALUES ")
	fieldArgs := make([]any, 0, len(orderFields)*len(i.values))
	for idx, value := range i.values {
		val := reflect.Indirect(reflect.ValueOf(value))
		if idx > 0 {
			i.sb.WriteString(", ")
		}
		i.sb.WriteByte('(')
		for count, field := range orderFields {
			if count > 0 {
				i.sb.WriteString(", ")
			}
			fd := val.FieldByName(field.FieldName)
			// 构建占位符 ？
			i.sb.WriteByte('?')
			// 存储字段数据
			fieldArgs = append(fieldArgs, fd.Interface())
		}
		i.sb.WriteByte(')')
	}
	i.args = append(i.args, fieldArgs...)
	return nil
}

// ExecuteWithContext 执行SQL语句
func (i *InsertSQL[T]) ExecuteWithContext(ctx context.Context) (*Result, error) {
	sqlInfo, err := i.Build()
	if err != nil {
		return nil, err
	}
	res, err := i.db.db.ExecContext(ctx, sqlInfo.SQL, sqlInfo.Args...)
	if err != nil {
		return &Result{
			err: err,
			res: nil,
		}, nil
	}
	return &Result{res: res}, err
}

// Build 构造SQL语句和维护SQL参数
// INSERT INTO `test_model` (`id`, `first_name`, `age`, `last_name`) VALUES (?, ?, ?, ?), (?, ?, ?, ?), (?, ?, ?, ?), (?, ?, ?, ?);
func (i *InsertSQL[T]) Build() (*SQLInfo, error) {
	// 构建SQL基本架构
	i.sb.WriteString("INSERT INTO ")
	var err error
	i.model, err = i.db.manager.Get(new(T))
	if err != nil {
		return nil, err
	}
	// 构建表名
	i.sb.WriteByte('`')
	i.sb.WriteString(i.model.TableName)
	i.sb.WriteByte('`')
	i.sb.WriteByte(' ')

	// TODO 构建COLUMNS 和 VALUES语句
	if err = i.buildColumnsAndValues(); err != nil {
		return nil, err
	}

	i.sb.WriteByte(';')
	res := &SQLInfo{SQL: i.sb.String(), Args: i.args}
	return res, nil
}

func NewInsertSQL[T any](db *DB) *InsertSQL[T] {
	return &InsertSQL[T]{
		// sb:   &strings.Builder{},
		// args: []any{},
		builder: newBuilder(),
		db:      db,
	}
}
