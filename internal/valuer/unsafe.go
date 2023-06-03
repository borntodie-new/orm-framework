package valuer

import (
	"database/sql"
	"github.com/borntodie-new/orm-framework/internal/errs"
	"github.com/borntodie-new/orm-framework/model"
	"reflect"
	"unsafe"
)

type unsafeValuer struct {
	// model 表模型，反射需要用到模型
	model *model.Model
	// addr 待解析的 T 结构体在内存中的地址
	addr unsafe.Pointer
}

func (u unsafeValuer) SetField(rows *sql.Rows) error {
	orderColumnsStr, err := rows.Columns()
	if err != nil {
		return err
	}
	receiptInterfaceFields := make([]any, 0, len(u.model.Fields))
	for _, str := range orderColumnsStr {
		fd, ok := u.model.ColumnsMap[str]
		if !ok {
			return errs.NewErrNotSupportUnknownColumn(str)
		}
		// 计算当前字段在 T 结构体中的偏移量
		address := unsafe.Pointer(uintptr(u.addr) + fd.Offset)
		receiptInterfaceFields = append(receiptInterfaceFields, reflect.NewAt(fd.Type, address).Interface())
	}
	return rows.Scan(receiptInterfaceFields...)
}

func (u unsafeValuer) GetField(fieldName string) (any, error) {
	fd, ok := u.model.FieldsMap[fieldName]
	if !ok {
		return nil, errs.NewErrNotSupportUnknownField(fieldName)
	}
	address := unsafe.Pointer(uintptr(u.addr) + fd.Offset)
	// 注意，在reflect中，不论是 New还是NewAt方法，返回的都是一个指针
	return reflect.NewAt(fd.Type, address).Elem().Interface(), nil
}

var _ FactoryValuer = NewUnsafeValuer

// NewUnsafeValuer entity必须是一个一级指针
func NewUnsafeValuer(model *model.Model, entity any) Valuer {
	return unsafeValuer{
		model: model,
		addr:  unsafe.Pointer(reflect.ValueOf(entity).Pointer()),
	}
}
