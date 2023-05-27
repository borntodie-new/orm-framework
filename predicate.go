package orm_framework

type opType string

func (o opType) String() string {
	return string(o)
}

const (
	EQType  = " = "
	GTType  = " > "
	GTEType = " >= "
	LTType  = " < "
	LTEType = " <= "
	ANDType = " AND "
	ORType  = " OR "
	NOTType = "NOT "
)

// Predicate 谓词，用于拼接WHERE条件的
// 需要实现 Expression 接口，因为 Predicate 可能会成为 left 或 right 字段数据
type Predicate struct {
	// field 字段名，注意这是Go中的字段名
	left Expression
	// op 操作符，具体是AND 或 OR 或 NOT 或 EQ 等等操作符
	op opType
	// value 字段值
	right Expression
}

// expr 标记位
func (p Predicate) expr() {
	//TODO implement me
	panic("implement me")
}

// AND 实现SQL中的 AND 语句
// Go中的使用：F("Id").EQ(12).AND(F("FirstName").EQ("Neo"))
// SQL中的使用：WHERE Id = 12 AND FirstName = "Neo"
func (p Predicate) AND(pre Predicate) Predicate {
	return Predicate{
		left:  p,
		op:    ANDType,
		right: pre,
	}
}

// OR 实现SQL中的 OR 语句
// Go中的使用：F("Id").EQ(12).OR(F("FirstName").EQ("Neo"))
// SQL中的使用：WHERE Id = 12 OR FirstName = "Neo"
func (p Predicate) OR(pre Predicate) Predicate {
	return Predicate{
		left:  p,
		op:    ORType,
		right: pre,
	}
}

// NOT 实现SQL中的 NOT 语句
// Go中的使用：NOT(F("Id").EQ(12))
// SQL中的使用：WHERE NOT Id = 12
// ⚠️：这里的pre应该是right
func NOT(pre Predicate) Predicate {
	return Predicate{
		op:    NOTType,
		right: pre,
	}
}
