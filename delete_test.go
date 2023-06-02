package orm_framework

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/borntodie-new/orm-framework/internal/errs"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestDeleteSQL_Build(t *testing.T) {
	db := memoryDB(t)
	testCases := []struct {
		name    string
		d       *DeleteSQL[TestModel]
		wantRes *SQLInfo
		wantErr error
	}{
		{
			name:    "test Table",
			d:       NewDeleteSQL[TestModel](db).Table("test_model"),
			wantRes: &SQLInfo{SQL: "DELETE FROM `test_model`;", Args: []any{}},
		},
		{
			name:    "test Where",
			d:       NewDeleteSQL[TestModel](db).Where(F("Id").EQ(12).AND(F("FirstName").EQ("Neo"))),
			wantRes: &SQLInfo{SQL: "DELETE FROM `test_model` WHERE (`id` = ?) AND (`first_name` = ?);", Args: []any{12, "Neo"}},
		},
		{
			name:    "test Table and Where",
			d:       NewDeleteSQL[TestModel](db).Table("test_model").Where(F("Id").EQ(12).AND(F("FirstName").EQ("Neo"))),
			wantRes: &SQLInfo{SQL: "DELETE FROM `test_model` WHERE (`id` = ?) AND (`first_name` = ?);", Args: []any{12, "Neo"}},
		},
		{
			name:    "test GT condition",
			d:       NewDeleteSQL[TestModel](db).Table("test_model").Where(F("Id").GT(12)),
			wantRes: &SQLInfo{SQL: "DELETE FROM `test_model` WHERE (`id` > ?);", Args: []any{12}},
		},
		{
			name:    "test GTE condition",
			d:       NewDeleteSQL[TestModel](db).Table("test_model").Where(F("Id").GTE(12)),
			wantRes: &SQLInfo{SQL: "DELETE FROM `test_model` WHERE (`id` >= ?);", Args: []any{12}},
		},
		{
			name:    "test LT condition",
			d:       NewDeleteSQL[TestModel](db).Table("test_model").Where(F("Id").LT(12)),
			wantRes: &SQLInfo{SQL: "DELETE FROM `test_model` WHERE (`id` < ?);", Args: []any{12}},
		},
		{
			name:    "test LTE condition",
			d:       NewDeleteSQL[TestModel](db).Table("test_model").Where(F("Id").LTE(12)),
			wantRes: &SQLInfo{SQL: "DELETE FROM `test_model` WHERE (`id` <= ?);", Args: []any{12}},
		},
		{
			name:    "test AND condition",
			d:       NewDeleteSQL[TestModel](db).Table("test_model").Where(F("Id").LTE(12).AND(F("FirstName").EQ("Neo"))),
			wantRes: &SQLInfo{SQL: "DELETE FROM `test_model` WHERE (`id` <= ?) AND (`first_name` = ?);", Args: []any{12, "Neo"}},
		},
		{
			name:    "test OR condition",
			d:       NewDeleteSQL[TestModel](db).Table("test_model").Where(F("Id").LTE(12).OR(F("FirstName").EQ("Neo"))),
			wantRes: &SQLInfo{SQL: "DELETE FROM `test_model` WHERE (`id` <= ?) OR (`first_name` = ?);", Args: []any{12, "Neo"}},
		},
		{
			name:    "test NOT condition",
			d:       NewDeleteSQL[TestModel](db).Table("test_model").Where(NOT(F("Id").EQ(12))),
			wantRes: &SQLInfo{SQL: "DELETE FROM `test_model` WHERE NOT (`id` = ?);", Args: []any{12}},
		},
		{
			name:    "test NOT and AND condition",
			d:       NewDeleteSQL[TestModel](db).Table("test_model").Where(NOT(F("Id").EQ(12)).AND(F("FirstName").EQ("Neo"))),
			wantRes: &SQLInfo{SQL: "DELETE FROM `test_model` WHERE NOT (`id` = ?) AND (`first_name` = ?);", Args: []any{12, "Neo"}},
		},
		{
			name:    "test not support unknown field",
			d:       NewDeleteSQL[TestModel](db).Where(F("id").EQ(12).AND(F("FirstName").EQ("Neo"))),
			wantErr: errs.NewErrNotSupportUnknownField("id"),
		},
		{
			name:    "test diy field name",
			d:       NewDeleteSQL[TestModel](db).Where(F("LastName").EQ("Jason")),
			wantRes: &SQLInfo{SQL: "DELETE FROM `test_model` WHERE (`test_model_last_name` = ?);", Args: []any{"Jason"}},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			res, err := tc.d.Build()
			assert.Equal(t, tc.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, tc.wantRes, res)
		})
	}
}

func TestDeleteSQL_ExecuteWithContext(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	mockDB, mock, err := sqlmock.New()
	assert.NoError(t, err)
	db, err := OpenDB(mockDB)

	testCases := []struct {
		name       string
		d          *DeleteSQL[TestModel]
		prepareSQL func()
		affected   int64
		wantErr    error
	}{
		{
			name: "no db",
			prepareSQL: func() {
				mock.ExpectExec("DELETE FROM .*").WillReturnError(errors.New("no db"))
			},
			d:       NewDeleteSQL[TestModel](db).Where(F("Id").EQ(12)),
			wantErr: errors.New("no db"),
		},
		{
			name: "affected success",
			prepareSQL: func() {
				result := driver.RowsAffected(19)
				mock.ExpectExec("DELETE FROM .*").WillReturnResult(result)
			},
			d:        NewDeleteSQL[TestModel](db).Where(F("Id").EQ(12)),
			affected: int64(19),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.prepareSQL()
			res, err := tc.d.ExecuteWithContext(ctx)
			assert.Equal(t, tc.wantErr, res.err)
			if err != nil {
				return
			}
			affected, err := res.RowsAffected()
			if err != nil {
				return
			}
			assert.Equal(t, tc.affected, affected)
		})
	}
}

type TestModel struct {
	Id        int8
	FirstName string
	Age       uint8
	LastName  *sql.NullString `orm:"column=test_model_last_name"`
}

// memoryDB 返回一个基于内存的 ORM，它使用的是 sqlite3 内存模式。
func memoryDB(t *testing.T) *DB {
	orm, err := Open("sqlite3", "file:test.db?cache=shared&mode=memory")
	if err != nil {
		t.Fatal(err)
	}
	return orm
}

// memoryDBWithDB 基于用户自定的库启动数据库
func memoryDBWithDB(db string, t *testing.T) *DB {
	orm, err := Open("sqlite3", fmt.Sprintf("file:%s.db?cache=shared&mode=memory", db))
	if err != nil {
		t.Fatal(err)
	}
	return orm
}
