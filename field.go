package orm_framework

type Field struct {
	fieldName string
}

// expr 完全是一个标记位，不做任何事情
func (f Field) expr() {
	//TODO implement me
	panic("implement me")
}

// F 初始化一个Field实例对象
func F(fieldName string) Field {
	return Field{fieldName: fieldName}
}

func (f Field) EQ(val any) Predicate {
	return Predicate{
		left:  f,
		op:    EQType,
		right: valueOf(val),
	}
}
func (f Field) GT(val any) Predicate {
	return Predicate{
		left:  f,
		op:    GTType,
		right: valueOf(val),
	}
}
func (f Field) GTE(val any) Predicate {
	return Predicate{
		left:  f,
		op:    GTEType,
		right: valueOf(val),
	}
}
func (f Field) LT(val any) Predicate {
	return Predicate{
		left:  f,
		op:    LTType,
		right: valueOf(val),
	}
}
func (f Field) LTE(val any) Predicate {
	return Predicate{
		left:  f,
		op:    LTEType,
		right: valueOf(val),
	}
}

/*
现在我们想想为什么上面5个方法需要返回 Predicate 实例对象
这要从我们的使用方法来说
F("Id").EQ(12) => Id > 12
F("Id").GTE(12) => Id >= 12

那这种使用场景怎么支持呢？
Id > 12 AND FirstName = "Neo"
从这个使用场景中，我们就能发现，现在是两个 Predicate 之间相互联系
所以，AND、OR 等语句需要在 Predicate层面上支持
*/
