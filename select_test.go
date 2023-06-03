package orm_framework

import (
	"context"
	"database/sql"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/borntodie-new/orm-framework/internal/errs"
	"github.com/borntodie-new/orm-framework/internal/valuer"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestSelectSQL_Build(t *testing.T) {
	db := memoryDB(t)
	testCases := []struct {
		name    string
		s       *SelectSQL[TestModel]
		wantRes *SQLInfo
		wantErr error
	}{
		{
			name: "test table",
			s:    NewSelectSQL[TestModel](db, valuer.NewUnsafeValuer),
			wantRes: &SQLInfo{
				SQL:  "SELECT * FROM `test_model`;",
				Args: []any{},
			},
		},
		{
			name: "test where",
			s:    NewSelectSQL[TestModel](db, valuer.NewUnsafeValuer).Where(F("Id").GTE(12)).Where(F("FirstName").EQ("JASON")),
			wantRes: &SQLInfo{
				SQL:  "SELECT * FROM `test_model` WHERE (`id` >= ?) AND (`first_name` = ?);",
				Args: []any{12, "JASON"},
			},
		},
		{
			name: "test specially fields",
			s:    NewSelectSQL[TestModel](db, valuer.NewUnsafeValuer).Fields(Common("Id"), Common("LastName")).Where(F("Id").GTE(12)).Where(F("FirstName").EQ("JASON")),
			wantRes: &SQLInfo{
				SQL:  "SELECT `id`, `test_model_last_name` FROM `test_model` WHERE (`id` >= ?) AND (`first_name` = ?);",
				Args: []any{12, "JASON"},
			},
		},
		{
			name:    "test with invalid specially fields",
			s:       NewSelectSQL[TestModel](db, valuer.NewUnsafeValuer).Fields(Common("Invalid")).Where(F("Id").GTE(12)).Where(F("FirstName").EQ("JASON")),
			wantErr: errs.NewErrNotSupportUnknownField("Invalid"),
		},
		{
			name:    "test with invalid where fields",
			s:       NewSelectSQL[TestModel](db, valuer.NewUnsafeValuer).Where(F("Invalid").GTE(12)),
			wantErr: errs.NewErrNotSupportUnknownField("Invalid"),
		},
		{
			name: "test AVG aggregate function",
			s:    NewSelectSQL[TestModel](db, valuer.NewUnsafeValuer).Fields(Avg("Age")).Where(F("Id").EQ(12)),
			wantRes: &SQLInfo{
				SQL:  "SELECT AVG(`age`) FROM `test_model` WHERE (`id` = ?);",
				Args: []any{12},
			},
		},
		{
			name: "test aggregate without where clause",
			s:    NewSelectSQL[TestModel](db, valuer.NewUnsafeValuer).Fields(Avg("Age")),
			wantRes: &SQLInfo{
				SQL:  "SELECT AVG(`age`) FROM `test_model`;",
				Args: []any{},
			},
		},
		{
			name: "test more aggregate function",
			s:    NewSelectSQL[TestModel](db, valuer.NewUnsafeValuer).Fields(Avg("Age"), Max("Age"), Count("Id")).Where(F("Id").EQ(12)),
			wantRes: &SQLInfo{
				SQL:  "SELECT AVG(`age`), MAX(`age`), COUNT(`id`) FROM `test_model` WHERE (`id` = ?);",
				Args: []any{12},
			},
		},
		{
			name:    "test invalid field name",
			s:       NewSelectSQL[TestModel](db, valuer.NewUnsafeValuer).Fields(Avg("Invalid")).Where(F("Id").EQ(12)),
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

func TestSelectSQL_QueryRawWithContext(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	mockDB, mock, err := sqlmock.New()
	assert.NoError(t, err)
	db, err := OpenDB(mockDB)

	testCases := []struct {
		name       string
		s          *SelectSQL[TestModel]
		prepareSQL func()
		wantRes    *TestModel
		wantErr    error
	}{
		{
			name: "test full columns",
			s:    NewSelectSQL[TestModel](db, valuer.NewUnsafeValuer).Where(F("Id").EQ(12)).Where(F("LastName").EQ("Neo")),
			prepareSQL: func() {
				mockRes := sqlmock.NewRows([]string{"id", "first_name", "age", "test_model_last_name"})
				mockRes.AddRow(12, "JASON", 18, "Neo")
				mock.ExpectQuery("SELECT .*").WillReturnRows(mockRes)
			},
			wantRes: &TestModel{
				Id:        12,
				FirstName: "JASON",
				Age:       18,
				LastName:  &sql.NullString{Valid: true, String: "Neo"},
			},
		},
		{
			name: "test specially columns",
			s:    NewSelectSQL[TestModel](db, valuer.NewUnsafeValuer).Fields(Common("Id"), Common("LastName")).Where(F("Id").EQ(12)).Where(F("LastName").EQ("Neo")),
			prepareSQL: func() {
				mockRes := sqlmock.NewRows([]string{"id", "test_model_last_name"})
				mockRes.AddRow(12, "Neo")
				mock.ExpectQuery("SELECT .*").WillReturnRows(mockRes)
			},
			wantRes: &TestModel{
				Id:       12,
				LastName: &sql.NullString{Valid: true, String: "Neo"},
			},
		},
		{
			name: "test invalid column",
			s:    NewSelectSQL[TestModel](db, valuer.NewUnsafeValuer).Fields(Common("Invalid")),
			prepareSQL: func() {
				mockRes := sqlmock.NewRows([]string{"id", "first_name", "age", "test_model_last_name"})
				mockRes.AddRow(12, "JASON", 18, "Neo")
				mock.ExpectQuery("SELECT .*").WillReturnRows(mockRes)
			},
			wantErr: errs.NewErrNotSupportUnknownField("Invalid"),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.prepareSQL()
			res, err := tc.s.QueryRawWithContext(ctx)
			assert.Equal(t, tc.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, tc.wantRes, res)
		})
	}
}

func TestSelectSQL_QueryWithContext(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	mockDB, mock, err := sqlmock.New()
	assert.NoError(t, err)
	db, err := OpenDB(mockDB)

	testCases := []struct {
		name       string
		s          *SelectSQL[TestModel]
		prepareSQL func()
		wantRes    []*TestModel
		wantErr    error
	}{
		{
			name: "test full columns",
			s:    NewSelectSQL[TestModel](db, valuer.NewUnsafeValuer).Where(F("Id").EQ(12)).Where(F("LastName").EQ("Neo")),
			prepareSQL: func() {
				mockRes := sqlmock.NewRows([]string{"id", "first_name", "age", "test_model_last_name"})
				mockRes.AddRow(12, "JASON", 18, "Neo")
				mockRes.AddRow(12, "JASON", 18, "Neo")
				mock.ExpectQuery("SELECT .*").WillReturnRows(mockRes)
			},
			wantRes: []*TestModel{
				{
					Id:        12,
					FirstName: "JASON",
					Age:       18,
					LastName:  &sql.NullString{Valid: true, String: "Neo"},
				},
				{
					Id:        12,
					FirstName: "JASON",
					Age:       18,
					LastName:  &sql.NullString{Valid: true, String: "Neo"},
				},
			},
		},
		{
			name: "test specially columns",
			s:    NewSelectSQL[TestModel](db, valuer.NewUnsafeValuer).Fields(Common("Id"), Common("LastName")).Where(F("Id").EQ(12)).Where(F("LastName").EQ("Neo")),
			prepareSQL: func() {
				mockRes := sqlmock.NewRows([]string{"id", "test_model_last_name"})
				mockRes.AddRow(12, "Neo")
				mockRes.AddRow(12, "Neo")
				mock.ExpectQuery("SELECT .*").WillReturnRows(mockRes)
			},
			wantRes: []*TestModel{
				{
					Id:       12,
					LastName: &sql.NullString{Valid: true, String: "Neo"},
				},
				{
					Id:       12,
					LastName: &sql.NullString{Valid: true, String: "Neo"},
				},
			},
		},
		{
			name: "test invalid column",
			s:    NewSelectSQL[TestModel](db, valuer.NewUnsafeValuer).Fields(Common("Invalid")),
			prepareSQL: func() {
				mockRes := sqlmock.NewRows([]string{"id", "first_name", "age", "test_model_last_name"})
				mockRes.AddRow(12, "JASON", 18, "Neo")
				mock.ExpectQuery("SELECT .*").WillReturnRows(mockRes)
			},
			wantErr: errs.NewErrNotSupportUnknownField("Invalid"),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.prepareSQL()
			res, err := tc.s.QueryWithContext(ctx)
			assert.Equal(t, tc.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, tc.wantRes, res)
		})
	}
}
