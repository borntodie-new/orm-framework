package errs

import (
	"errors"
	"fmt"
)

var (
	ErrNotSupportPredicate = errors.New("不支持Predicate类型")
	ErrNotSupportModelType = errors.New("不支持模型类型")
)

func NewErrNotSupportUnknownField(val any) error {
	return errors.New(fmt.Sprintf("不支持未知字段 %v ", val))
}
