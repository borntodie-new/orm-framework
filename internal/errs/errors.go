package errs

import (
	"errors"
	"fmt"
)

var (
	ErrNotSupportPredicate      = errors.New("不支持Predicate类型")
	ErrNotSupportModelType      = errors.New("不支持模型类型")
	ErrNotUpdateSQLSetClause    = errors.New("更新语句没有SET子句")
	ErrNotInsertSQLValuesClause = errors.New("插入语句没有VALUES子句")
)

func NewErrNotSupportUnknownField(val any) error {
	return errors.New(fmt.Sprintf("不支持未知字段 %v ", val))
}
func NewErrInvalidTagContext(val string) error {
	return errors.New(fmt.Sprintf("不支持标签文本 %s ", val))
}
