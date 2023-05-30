package orm_framework

import (
	"github.com/borntodie-new/orm-framework/internal/errs"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSelectSQL_Build(t *testing.T) {
	db := memoryDB(t)
	testCases := []struct {
		name    string
		s       *SelectSQL[*TestModel]
		wantRes *SQLInfo
		wantErr error
	}{
		{
			name: "test table",
			s:    NewSelectSQL[*TestModel](db).Table("db_test_model"),
			wantRes: &SQLInfo{
				SQL:  "SELECT `id`, `first_name`, `age`, `test_model_last_name` FROM `db_test_model`;",
				Args: []any{},
			},
		},
		{
			name: "test where",
			s:    NewSelectSQL[*TestModel](db).Where(F("Id").GTE(12)).Where(F("FirstName").EQ("JASON")),
			wantRes: &SQLInfo{
				SQL:  "SELECT `id`, `first_name`, `age`, `test_model_last_name` FROM `test_model` WHERE (`id` >= ?) AND (`first_name` = ?);",
				Args: []any{12, "JASON"},
			},
		},
		{
			name: "test specially fields",
			s:    NewSelectSQL[*TestModel](db).Fields("Id", "LastName").Where(F("Id").GTE(12)).Where(F("FirstName").EQ("JASON")),
			wantRes: &SQLInfo{
				SQL:  "SELECT `id`, `test_model_last_name` FROM `test_model` WHERE (`id` >= ?) AND (`first_name` = ?);",
				Args: []any{12, "JASON"},
			},
		},
		{
			name:    "test with invalid specially fields",
			s:       NewSelectSQL[*TestModel](db).Fields("Invalid").Where(F("Id").GTE(12)).Where(F("FirstName").EQ("JASON")),
			wantErr: errs.NewErrNotSupportUnknownField("Invalid"),
		},
		{
			name:    "test with invalid where fields",
			s:       NewSelectSQL[*TestModel](db).Where(F("Invalid").GTE(12)),
			wantErr: errs.NewErrNotSupportUnknownField("Invalid"),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			res, err := tc.s.Build()
			assert.Equal(t, tc.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, tc.wantRes, res)
		})
	}

}
