package model

import (
	"github.com/borntodie-new/orm-framework/internal/errs"
	"reflect"
	"sync"
	"unicode"
)

// 统一管理model表模型的结构

// Manager 统一管理model表模型结构
type Manager struct {
	// models 需要管理的所有model模型
	// 为什么用 sync.Map 结构呢？因为这个可以避免并发问题
	models sync.Map
}

// Get 获取表模型
// key参数其实就是模型
// 我们这里是用模型在Go中的Type类型作为key
func (m *Manager) Get(key any) (*Model, error) {
	typ := reflect.TypeOf(key)
	if typ.Kind() != reflect.Pointer || typ.Elem().Kind() != reflect.Struct {
		return nil, errs.ErrNotSupportModelType
	}
	// 因为model是一个一级指针类型的结构体，需要将这个指针结构体的本体找到
	// 因此需要使用 Elem() 方法
	typ = typ.Elem()
	mod, ok := m.models.Load(typ)
	if ok {
		return mod.(*Model), nil
	}
	return m.Register(typ)
}

// Register 注册表模型
func (m *Manager) Register(typ reflect.Type) (*Model, error) {
	// 构建数据
	fieldsMap := make(map[string]*Field)
	columnsMap := make(map[string]*Field)
	for i := 0; i < typ.NumField(); i++ {
		fd := typ.Field(i)
		f := &Field{
			FieldName:  fd.Name,
			ColumnName: underscoreName(fd.Name),
			Type:       fd.Type,
		}
		fieldsMap[fd.Name] = f
		columnsMap[underscoreName(fd.Name)] = f
	}
	mod := &Model{
		TableName:  underscoreName(typ.Name()),
		Type:       typ,
		FieldsMap:  fieldsMap,
		ColumnsMap: columnsMap,
	}
	m.models.Store(typ, mod)
	return mod, nil
}

// underscoreName 驼峰转字符串命名
func underscoreName(tableName string) string {
	var buf []byte
	for i, v := range tableName {
		if unicode.IsUpper(v) {
			if i != 0 {
				buf = append(buf, '_')
			}
			buf = append(buf, byte(unicode.ToLower(v)))
		} else {
			buf = append(buf, byte(v))
		}

	}
	return string(buf)
}
