package orm_framework

import "database/sql"

// Result ExecuteSQL 统一返回的结果信息
type Result struct {
	// err 执行SQL出现的错误信息
	err error
	// res 执行SQL返回的结果
	res sql.Result
}

func (r *Result) LastInsertId() (int64, error) {
	if r.err != nil {
		return 0, r.err
	}
	return r.res.LastInsertId()
}

func (r *Result) RowsAffected() (int64, error) {
	if r.err != nil {
		return 0, r.err
	}
	return r.res.RowsAffected()
}
