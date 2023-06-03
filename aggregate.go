package orm_framework

// 聚合函数
// 看下聚合函数怎么用
// SELECT AVG(`id`), COUNT(`id`) FROM `test_model`;

type AggregateFn string

func (a AggregateFn) String() string {
	return string(a)
}

const (
	AVG   AggregateFn = "AVG"
	MAX   AggregateFn = "MAX"
	MIN   AggregateFn = "MIN"
	COUNT AggregateFn = "COUNT"
	SUM   AggregateFn = "SUM"
)

type Aggregate struct {
	// fieldName Go中结构体的字段名
	fieldName string
	// fn 聚合函数名
	fn AggregateFn
}

// Avg 平均数
func Avg(fieldName string) Aggregate {
	return Aggregate{
		fieldName: fieldName,
		fn:        AVG,
	}
}

// Max 最大值
func Max(fieldName string) Aggregate {
	return Aggregate{
		fieldName: fieldName,
		fn:        MAX,
	}
}

// Min 最小值
func Min(fieldName string) Aggregate {
	return Aggregate{
		fieldName: fieldName,
		fn:        MIN,
	}
}

// Count 数据个数
func Count(fieldName string) Aggregate {
	return Aggregate{
		fieldName: fieldName,
		fn:        COUNT,
	}
}

// Sum 数据和
func Sum(fieldName string) Aggregate {
	return Aggregate{
		fieldName: fieldName,
		fn:        SUM,
	}
}

func Common(fieldName string) Aggregate {
	return Aggregate{fieldName: fieldName}
}
