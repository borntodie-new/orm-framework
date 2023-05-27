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
			wantRes: &SQLInfo{SQL: "DELETE FROM `test_model`;"},
		},
		{
			name:    "test Where",
			d:       NewDeleteSQL[TestModel]().Where("id").Where("first_name"),
			wantRes: &SQLInfo{SQL: "DELETE FROM `TestModel` WHERE (`id` = ?) AND (`first_name` = ?);"},
		},
		{
			name:    "test Table and Where",
			d:       NewDeleteSQL[TestModel]().Table("test_model").Where("id").Where("first_name"),
			wantRes: &SQLInfo{SQL: "DELETE FROM `test_model` WHERE (`id` = ?) AND (`first_name` = ?);"},
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
