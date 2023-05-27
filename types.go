package orm_framework

import (
	"context"
	"database/sql"
)

// Querier 查询语句接口，主要是用于SELECT语句
// 功能是执行Query和QueryRaw
// 这里使用了范型，为了约束类型的
type Querier[T any] interface {
	QueryWithContext(ctx context.Context) (*T, error)
	QueryRawWithContext(ctx context.Context) ([]*T, error)
}

// Executer 执行语句接口，主要是用于DELETE、UPDATE、UPDATE语句
// 功能是执行SQL语句
// 注意，这里不需要使用范型，因为我们并不需要将结果返回
type Executer interface {
	ExecuteWithContext(ctx context.Context) (sql.Result, error)
}

// Builder 构建SQL语句的的接口
type Builder interface {
	// Build 构建SQL语句
	Build() (*SQLInfo, error)
}

// SQLInfo 保存SQL语句的SQL参数的结构体
type SQLInfo struct {
	// SQL 具体的SQL语句
	SQL string
	// Args 具体的SQL参数
	Args any
}
