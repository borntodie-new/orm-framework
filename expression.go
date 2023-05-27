package orm_framework

// Expression 表达式，其实就是一个标记位，用于约束表达式的类型的
type Expression interface {
	expr()
}

// Value 作为Predicate 的right
// 但是Predicate的right是一个 Expression 类型，所以我们也需要将 Value 设置成 Expression 类型
type Value struct {
	val any
}

// expr 完全是一个标记位，不做任何处理
func (v Value) expr() {
	//TODO implement me
	panic("implement me")
}

func valueOf(val any) Value {
	return Value{val: val}
}
