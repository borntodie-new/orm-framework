package valuer

import (
	"database/sql"
	"github.com/borntodie-new/orm-framework/internal/errs"
	"github.com/borntodie-new/orm-framework/model"
	"reflect"
)

type reflectValuer struct {
	// model 表模型，反射需要用到模型
	model *model.Model
	// t 待解析的 T 结构体
	t reflect.Value
}

func (r reflectValuer) SetField(rows *sql.Rows) error {
	orderColumnsStr, err := rows.Columns()
	if err != nil {
		return err
	}
	receiptFields := make([]any, 0, len(r.model.Fields))                    // 保存Scan需要数据
	receiptInterfaceFields := make([]reflect.Value, 0, len(r.model.Fields)) // 用于保存每个字段的 Value类型
	for _, str := range orderColumnsStr {
		fd, ok := r.model.ColumnsMap[str]
		if !ok {
			return errs.NewErrNotSupportUnknownColumn(str)
		}
		temp := reflect.New(fd.Type)
		receiptFields = append(receiptFields, temp.Interface())
		receiptInterfaceFields = append(receiptInterfaceFields, temp.Elem())
	}
	// 接收SQL返回的结果数据
	err = rows.Scan(receiptFields...)
	if err != nil {
		return err
	}
	// 将Scan出来的数据设置到 tp 结构体字段上
	for idx, str := range orderColumnsStr {
		fd, ok := r.model.ColumnsMap[str]
		if !ok {
			return errs.NewErrNotSupportUnknownColumn(str)
		}
		r.t.FieldByName(fd.FieldName).Set(receiptInterfaceFields[idx])
	}
	return nil
}

func (r reflectValuer) GetField(fieldName string) (any, error) {
	fd, ok := r.model.FieldsMap[fieldName]
	if !ok {
		return nil, errs.NewErrNotSupportUnknownField(fieldName)
	}
	return r.t.FieldByName(fd.FieldName).Interface(), nil
}

var _ FactoryValuer = NewReflectValuer

// NewReflectValuer entity必须是一个一级指针
func NewReflectValuer(model *model.Model, entity any) Valuer {
	return reflectValuer{
		model: model,
		t:     reflect.Indirect(reflect.ValueOf(entity)),
	}
}
