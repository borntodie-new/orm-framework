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

func TestRawSQL_Build(t *testing.T) {
	db := memoryDB(t)
	testCases := []struct {
		name    string
		r       *RawSQL[TestModel]
		wantRes *SQLInfo
		wantErr error
	}{
		{
			name:    "test sample",
			r:       NewRawSQL[TestModel](db, valuer.NewReflectValuer, "SELECT * FROM `test_model`;"),
			wantRes: &SQLInfo{SQL: "SELECT * FROM `test_model`;"},
		},
		{
			name:    "test sample with args",
			r:       NewRawSQL[TestModel](db, valuer.NewUnsafeValuer, "SELECT * FROM `test_model` WHERE `id` = ?;", 12),
			wantRes: &SQLInfo{SQL: "SELECT * FROM `test_model` WHERE `id` = ?;", Args: []any{12}},
		},
		{
			name:    "test no SQL",
			r:       NewRawSQL[TestModel](db, valuer.NewUnsafeValuer, ""),
			wantErr: errs.ErrNoSQL,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			res, err := tc.r.Build()
			assert.Equal(t, tc.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, tc.wantRes, res)
		})
	}
}

func TestRawSQL_QueryRawWithContext(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	mockDB, mock, err := sqlmock.New()
	assert.NoError(t, err)
	db, err := OpenDB(mockDB)

	testCases := []struct {
		name       string
		s          *RawSQL[TestModel]
		prepareSQL func()
		wantRes    *TestModel
		wantErr    error
	}{
		{
			name: "test full columns",
			s:    NewRawSQL[TestModel](db, valuer.NewUnsafeValuer, "SELECT `id`, `first_name`, `age` `test_model_last_name` FROM `test_model`;"),
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
			s:    NewRawSQL[TestModel](db, valuer.NewReflectValuer, "SELECT `id`, `test_model_last_name` FROM `test_model`;"),
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
			name: "test no sql",
			s:    NewRawSQL[TestModel](db, valuer.NewUnsafeValuer, ""),
			prepareSQL: func() {
				mockRes := sqlmock.NewRows([]string{"id", "first_name", "age", "test_model_last_name"})
				mockRes.AddRow(12, "JASON", 18, "Neo")
				mock.ExpectQuery("SELECT .*").WillReturnRows(mockRes)
			},
			wantErr: errs.ErrNoSQL,
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

func TestRawSQL_QueryWithContext(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	mockDB, mock, err := sqlmock.New()
	assert.NoError(t, err)
	db, err := OpenDB(mockDB)

	testCases := []struct {
		name       string
		s          *RawSQL[TestModel]
		prepareSQL func()
		wantRes    []*TestModel
		wantErr    error
	}{
		{
			name: "test full columns",
			s:    NewRawSQL[TestModel](db, valuer.NewUnsafeValuer, "SELECT `id`, `first_name`, `age` `test_model_last_name` FROM `test_model`;"),
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
				}, {
					Id:        12,
					FirstName: "JASON",
					Age:       18,
					LastName:  &sql.NullString{Valid: true, String: "Neo"},
				},
			},
		},
		{
			name: "test specially columns",
			s:    NewRawSQL[TestModel](db, valuer.NewReflectValuer, "SELECT `id`, `test_model_last_name` FROM `test_model`;"),
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
				}, {
					Id:       12,
					LastName: &sql.NullString{Valid: true, String: "Neo"},
				},
			},
		},
		{
			name: "test no sql",
			s:    NewRawSQL[TestModel](db, valuer.NewUnsafeValuer, ""),
			prepareSQL: func() {
				mockRes := sqlmock.NewRows([]string{"id", "first_name", "age", "test_model_last_name"})
				mockRes.AddRow(12, "JASON", 18, "Neo")
				mock.ExpectQuery("SELECT .*").WillReturnRows(mockRes)
			},
			wantErr: errs.ErrNoSQL,
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
