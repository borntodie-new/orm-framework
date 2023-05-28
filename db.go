package orm_framework

import (
	"database/sql"
	"github.com/borntodie-new/orm-framework/model"
)

type DB struct {
	// db 真实客SQL做交互的数据库连接对象
	db *sql.DB
	// manager model 管理器
	manager *model.Manager
}

// Open 创建自定义的 DB 实例对象
func Open(driver string, dataSourceName string) (*DB, error) {
	db, err := sql.Open(driver, dataSourceName)
	if err != nil {
		return nil, err
	}
	return OpenDB(db)
}

// OpenDB 创建自定义的 DB 实例对象
// 疑问：为什么已经有了 Open 方法，还需要提供这个方法
// 为了扩展性，这也是 Go 内置的 sql 的设计传统
func OpenDB(db *sql.DB) (*DB, error) {
	return &DB{db: db, manager: &model.Manager{}}, nil
}
