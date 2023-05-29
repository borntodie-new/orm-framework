package orm_framework

import (
	"database/sql"
	"github.com/borntodie-new/orm-framework/internal/errs"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestInsertSQL_Build(t *testing.T) {
	db := memoryDB(t)
	testCases := []struct {
		name    string
		i       *InsertSQL[*TestModel]
		wantRes *SQLInfo
		wantErr error
	}{
		{
			name: "test one values with full fields",
			i: NewInsertSQL[*TestModel](db).Values(&TestModel{
				Id:        1,
				FirstName: "Jason",
				Age:       19,
				LastName:  &sql.NullString{Valid: true, String: "Neo"},
			}),
			wantRes: &SQLInfo{
				SQL:  "INSERT INTO `test_model` (`id`, `first_name`, `age`, `test_model_last_name`) VALUES (?, ?, ?, ?);",
				Args: []any{int8(1), "Jason", uint8(19), &sql.NullString{Valid: true, String: "Neo"}},
			},
		},
		{
			name: "test multiple values with full fields",
			i: NewInsertSQL[*TestModel](db).Values(&TestModel{
				Id:        1,
				FirstName: "Jason",
				Age:       19,
				LastName:  &sql.NullString{Valid: true, String: "Neo"},
			}, &TestModel{
				Id:        100,
				FirstName: "Tank",
				Age:       67,
				LastName:  &sql.NullString{Valid: true, String: "Alice"},
			}),
			wantRes: &SQLInfo{
				SQL: "INSERT INTO `test_model` (`id`, `first_name`, `age`, `test_model_last_name`) VALUES (?, ?, ?, ?), (?, ?, ?, ?);",
				Args: []any{
					int8(1), "Jason", uint8(19), &sql.NullString{Valid: true, String: "Neo"},
					int8(100), "Tank", uint8(67), &sql.NullString{Valid: true, String: "Alice"},
				},
			},
		},
		{
			name: "test one values with specially fields",
			i: NewInsertSQL[*TestModel](db).Fields("Id", "LastName").Values(&TestModel{
				Id:       1,
				LastName: &sql.NullString{Valid: true, String: "Neo"},
			}),
			wantRes: &SQLInfo{
				SQL:  "INSERT INTO `test_model` (`id`, `test_model_last_name`) VALUES (?, ?);",
				Args: []any{int8(1), &sql.NullString{Valid: true, String: "Neo"}},
			},
		},
		{
			name: "test multiple values with specially fields",
			i: NewInsertSQL[*TestModel](db).Fields("Id", "LastName").Values(&TestModel{
				Id:       1,
				LastName: &sql.NullString{Valid: true, String: "Neo"},
			}, &TestModel{
				Id:       100,
				LastName: &sql.NullString{Valid: true, String: "JASON"},
			}),
			wantRes: &SQLInfo{
				SQL: "INSERT INTO `test_model` (`id`, `test_model_last_name`) VALUES (?, ?), (?, ?);",
				Args: []any{
					int8(1), &sql.NullString{Valid: true, String: "Neo"},
					int8(100), &sql.NullString{Valid: true, String: "JASON"},
				},
			},
		},
		{
			name: "test multiple values with specially fields",
			i: NewInsertSQL[*TestModel](db).Fields("Invalid").Values(&TestModel{
				Id:       1,
				LastName: &sql.NullString{Valid: true, String: "Neo"},
			}),
			wantErr: errs.NewErrNotSupportUnknownField("Invalid"),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			res, err := tc.i.Build()
			assert.Equal(t, tc.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, tc.wantRes, res)
		})
	}
}
