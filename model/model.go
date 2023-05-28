package model

import "reflect"

const (
	FieldTagName  = "orm"
	ColumnTagName = "column"
)

// 存储表模型

// Model 表模型元数据
type Model struct {
	// tableName 表模型名
	TableName string
	// Type Go中结构体的类型
	Type reflect.Type
	// FieldsMap 保存表模型字段信息
	// Go中的字段名作为key
	FieldsMap map[string]*Field
	// ColumnsMap 保存表模型字段信息
	// SQL中的列名作为key
	ColumnsMap map[string]*Field
}

// Field Go中字段元数据
type Field struct {
	// Go中的字段名
	FieldName string
	// SQL中的列名
	ColumnName string
	// Type 字段在Go中的类型
	Type reflect.Type
}
