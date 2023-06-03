package orm_framework

import (
	"github.com/borntodie-new/orm-framework/model"
	"strings"
)

// builder 公用的构建 SQL 语句的结构
// 内部维护的都是涉及到构建 SQL 语句和 SQL 参数的数据
type builder struct {
	// sb 构建 SQL 语句的，性能好
	sb *strings.Builder
	// args SQL 语句中的参数
	args []any
	// model 表模型
	model *model.Model
}

func newBuilder() *builder {
	return &builder{
		sb:   &strings.Builder{},
		args: []any{},
	}
}
