package orm_framework

import (
	"github.com/borntodie-new/orm-framework/internal/errs"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUpdateSQL_Build(t *testing.T) {
	db := memoryDB(t)
	testCases := []struct {
		name    string
		u       *UpdateSQL[*TestModel]
		wantRes *SQLInfo
		wantErr error
	}{
		{
			name: "test table name",
			u:    NewUpdateSQL[*TestModel](db).Table("Order_TestModel").Values("Id", 1),
			wantRes: &SQLInfo{
				SQL:  "UPDATE `Order_TestModel` SET `id` = ?;",
				Args: []any{1},
			},
		},
		{
			name: "test set one clause",
			u:    NewUpdateSQL[*TestModel](db).Values("Id", 1),
			wantRes: &SQLInfo{
				SQL:  "UPDATE `test_model` SET `id` = ?;",
				Args: []any{1},
			},
		},
		{
			name: "test set multiple clause",
			u:    NewUpdateSQL[*TestModel](db).Values("Id", 1).Values("FirstName", "Neo"),
			wantRes: &SQLInfo{
				SQL:  "UPDATE `test_model` SET `id` = ?, `first_name` = ?;",
				Args: []any{1, "Neo"},
			},
		},
		{
			name: "test update with where",
			u:    NewUpdateSQL[*TestModel](db).Values("Id", 1).Values("FirstName", "Neo").Where(F("Age").GTE(12)),
			wantRes: &SQLInfo{
				SQL:  "UPDATE `test_model` SET `id` = ?, `first_name` = ? WHERE (`age` >= ?);",
				Args: []any{1, "Neo", 12},
			},
		},
		{
			name:    "test unknown field",
			u:       NewUpdateSQL[*TestModel](db).Values("Invalid", 1),
			wantErr: errs.NewErrNotSupportUnknownField("Invalid"),
		},
		{
			name:    "test no set clause",
			u:       NewUpdateSQL[*TestModel](db),
			wantErr: errs.ErrNotUpdateSQLSetClause,
		},
		{
			name:    "test no set clause with where",
			u:       NewUpdateSQL[*TestModel](db).Where(F("Id").EQ(12)),
			wantErr: errs.ErrNotUpdateSQLSetClause,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			res, err := tc.u.Build()
			assert.Equal(t, tc.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, tc.wantRes, res)
		})
	}
}