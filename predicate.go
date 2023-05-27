package orm_framework

type opType string

type Predicate struct {
	// field 字段名，注意这是Go中的字段名
	field string
	// op 操作符，具体是AND 或 OR 或 NOT 或 EQ 等等操作符
	op string
	// value 字段值
	value any
}

// P 实例化一个条件，使用链式调用的方式
func P(fieldName string) Predicate {
	return Predicate{field: fieldName}
}

// EQ SQL 中的 = 操作
func (p Predicate) EQ(value any) Predicate {
	p.op = "="
	p.value = value
	return p
}
