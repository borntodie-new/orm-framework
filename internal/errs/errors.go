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
	ErrNoRows                   = errors.New("没有查询到数据")
	ErrUnsupportedNil           = errors.New("不支持空指针类型")
	ErrNoSQL                    = errors.New("SQL语句不能为空")
)

func NewErrNotSupportUnknownField(val any) error {
	return errors.New(fmt.Sprintf("不支持未知字段 %v ", val))
}
func NewErrInvalidTagContext(val string) error {
	return errors.New(fmt.Sprintf("不支持标签文本 %s ", val))
}

func NewErrNotSupportUnknownColumn(val any) error {
	return errors.New(fmt.Sprintf("不支持未知列名 %v ", val))
}
