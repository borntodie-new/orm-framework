package valuer

import (
	"database/sql"
	"github.com/borntodie-new/orm-framework/model"
)

type Valuer interface {
	// SetField 将 SQL 中的数据映射到 Go 的结构体字段中
	SetField(rows *sql.Rows) error
	// GetField 将 Go 中的结构体上的字段数据返回
	GetField(fieldName string) (any, error)
}

// FactoryValuer 一个简单的工厂，用于返回 Valuer 类型的实现
type FactoryValuer func(model *model.Model, entity any) Valuer
