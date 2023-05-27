package orm_framework

import (
	"database/sql"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDeleteSQL_Build(t *testing.T) {
	testCases := []struct {
		name    string
		d       *DeleteSQL[TestModel]
		wantRes *SQLInfo
		wantErr error
	}{
		{
			name:    "test Table",
			d:       NewDeleteSQL[TestModel]().Table("test_model"),
			wantRes: &SQLInfo{SQL: "DELETE FROM `test_model`;", Args: []any{}},
		},
		{
			name:    "test Where",
			d:       NewDeleteSQL[TestModel]().Where(F("id").EQ(12).AND(F("first_name").EQ("Neo"))),
			wantRes: &SQLInfo{SQL: "DELETE FROM `TestModel` WHERE (`id` = ?) AND (`first_name` = ?);", Args: []any{12, "Neo"}},
		},
		{
			name:    "test Table and Where",
			d:       NewDeleteSQL[TestModel]().Table("test_model").Where(F("id").EQ(12).AND(F("first_name").EQ("Neo"))),
			wantRes: &SQLInfo{SQL: "DELETE FROM `test_model` WHERE (`id` = ?) AND (`first_name` = ?);", Args: []any{12, "Neo"}},
		},
		{
			name:    "test GT condition",
			d:       NewDeleteSQL[TestModel]().Table("test_model").Where(F("id").GT(12)),
			wantRes: &SQLInfo{SQL: "DELETE FROM `test_model` WHERE (`id` > ?);", Args: []any{12}},
		},
		{
			name:    "test GTE condition",
			d:       NewDeleteSQL[TestModel]().Table("test_model").Where(F("id").GTE(12)),
			wantRes: &SQLInfo{SQL: "DELETE FROM `test_model` WHERE (`id` >= ?);", Args: []any{12}},
		},
		{
			name:    "test LT condition",
			d:       NewDeleteSQL[TestModel]().Table("test_model").Where(F("id").LT(12)),
			wantRes: &SQLInfo{SQL: "DELETE FROM `test_model` WHERE (`id` < ?);", Args: []any{12}},
		},
		{
			name:    "test LTE condition",
			d:       NewDeleteSQL[TestModel]().Table("test_model").Where(F("id").LTE(12)),
			wantRes: &SQLInfo{SQL: "DELETE FROM `test_model` WHERE (`id` <= ?);", Args: []any{12}},
		},
		{
			name:    "test AND condition",
			d:       NewDeleteSQL[TestModel]().Table("test_model").Where(F("id").LTE(12).AND(F("first_name").EQ("Neo"))),
			wantRes: &SQLInfo{SQL: "DELETE FROM `test_model` WHERE (`id` <= ?) AND (`first_name` = ?);", Args: []any{12, "Neo"}},
		},
		{
			name:    "test OR condition",
			d:       NewDeleteSQL[TestModel]().Table("test_model").Where(F("id").LTE(12).OR(F("first_name").EQ("Neo"))),
			wantRes: &SQLInfo{SQL: "DELETE FROM `test_model` WHERE (`id` <= ?) OR (`first_name` = ?);", Args: []any{12, "Neo"}},
		},
		{
			name:    "test NOT condition",
			d:       NewDeleteSQL[TestModel]().Table("test_model").Where(NOT(F("id").EQ(12))),
			wantRes: &SQLInfo{SQL: "DELETE FROM `test_model` WHERE NOT (`id` = ?);", Args: []any{12}},
		},
		{
			name:    "test NOT and AND condition",
			d:       NewDeleteSQL[TestModel]().Table("test_model").Where(NOT(F("id").EQ(12)).AND(F("first_name").EQ("Neo"))),
			wantRes: &SQLInfo{SQL: "DELETE FROM `test_model` WHERE NOT (`id` = ?) AND (`first_name` = ?);", Args: []any{12, "Neo"}},
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

type TestModel struct {
	Id        int8
	FirstName string
	Age       uint8
	LastName  *sql.NullString
}
