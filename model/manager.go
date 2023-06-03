package model

import (
	"github.com/borntodie-new/orm-framework/internal/errs"
	"reflect"
	"strings"
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
	if key == nil {
		return nil, errs.ErrUnsupportedNil
	}
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
	return m.register(key)
}

func (m *Manager) register(key any) (*Model, error) {
	// 因为在 Get 方法中已经做了判断，所以这里直接用就好
	typ := reflect.TypeOf(key).Elem()
	// 构建数据
	numField := typ.NumField()
	fieldsMap := make(map[string]*Field, numField)
	columnsMap := make(map[string]*Field, numField)
	fields := make([]*Field, 0, numField)
	for i := 0; i < numField; i++ {
		fd := typ.Field(i)
		tagsMap, err := m.parseTag(fd.Tag)
		if err != nil {
			return nil, err
		}
		f := &Field{
			FieldName: fd.Name,
			Type:      fd.Type,
			Offset:    fd.Offset,
		}
		colName, ok := tagsMap[ColumnTagName]
		if ok && colName != "" {
			f.ColumnName = colName
		} else {
			f.ColumnName = underscoreName(fd.Name)
		}

		fieldsMap[fd.Name] = f
		// columnsMap[underscoreName(fd.Name)] = f
		columnsMap[f.ColumnName] = f
		fields = append(fields, f)
	}
	// 注意：这里的 TableName 接口不能定义在 ORM 框架的那个包中，因为会出现 循环引入 的问题
	var tableName string
	tbn, ok := key.(TableName)
	if ok {
		tableName = tbn.TableName()
	}
	if tableName == "" {
		tableName = underscoreName(typ.Name())
	}
	mod := &Model{
		TableName:  tableName,
		FieldsMap:  fieldsMap,
		ColumnsMap: columnsMap,
		Fields:     fields,
	}
	m.models.Store(typ, mod)
	return mod, nil
}

// Register 注册表模型
// 这个方法接收的参数是一个 reflect.Type 类型，他只能由 Get 方法调用
// 由于在 Get 方法内部需要对 T 进行反射，这里也需要反射，所以我们才设计成接收一个 reflect.Type 类型的参数
// 也是处于性能的考虑
// 现在遇到的问题是，我们想通过 interface 的形式对表模型数据进行设置模型名字
// 这就需要我们在这个方法中必须要有 T
//func (m *Manager) register(typ reflect.Type) (*Model, error) {
//	// 构建数据
//	numField := typ.NumField()
//	fieldsMap := make(map[string]*Field, numField)
//	columnsMap := make(map[string]*Field, numField)
//	fields := make([]*Field, 0, numField)
//	for i := 0; i < numField; i++ {
//		fd := typ.Field(i)
//		tagsMap, err := m.parseTag(fd.Tag)
//		if err != nil {
//			return nil, err
//		}
//		f := &Field{
//			FieldName: fd.Name,
//			Type:      fd.Type,
//		}
//		colName, ok := tagsMap[ColumnTagName]
//		if ok && colName != "" {
//			f.ColumnName = colName
//		} else {
//			f.ColumnName = underscoreName(fd.Name)
//		}
//
//		fieldsMap[fd.Name] = f
//		// columnsMap[underscoreName(fd.Name)] = f
//		columnsMap[f.ColumnName] = f
//		fields = append(fields, f)
//	}
//	mod := &Model{
//		TableName: underscoreName(typ.Name()),
//		//Type:       typ,
//		FieldsMap:  fieldsMap,
//		ColumnsMap: columnsMap,
//		Fields:     fields,
//	}
//	m.models.Store(typ, mod)
//	return mod, nil
//}

func (m *Manager) parseTag(tag reflect.StructTag) (map[string]string, error) {
	res := make(map[string]string, 1)
	tagStr, ok := tag.Lookup(FieldTagName)
	if !ok {
		return res, nil
	}
	pairs := strings.Split(tagStr, ",")
	for _, pair := range pairs {
		temp := strings.Split(pair, "=")
		if len(temp) != 2 {
			return nil, errs.NewErrInvalidTagContext(pair)
		}
		res[temp[0]] = temp[1]
	}
	return res, nil
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
