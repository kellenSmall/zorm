package demo

import (
	"context"
	"database/sql"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"zorm/demo/internal/errs"
)

func memoryDB(t *testing.T) *DB {
	orm, err := Open("sqlite3", "file:test.db?cache=shared&mode=memory")
	if err != nil {
		t.Fatal(err)
	}
	return orm
}
func TestSelector_Build(t *testing.T) {
	db := memoryDB(t)
	tests := []struct {
		name    string
		s       QueryBuilder
		want    *Query
		wantErr error
	}{
		{
			name: "from",
			s:    NewSelector[TestModel](db).From("test_model_tab"),
			want: &Query{
				SQL: "SELECT * FROM test_model_tab;",
			},
		},
		{
			name: "no from",
			s:    NewSelector[TestModel](db),
			want: &Query{
				SQL: "SELECT * FROM `TestModel`;",
			},
		},
		{
			name: "empty from",
			s:    NewSelector[TestModel](db).From(""),
			want: &Query{
				SQL: "SELECT * FROM `TestModel`;",
			},
		}, {
			name: "sigle and simple predicate",
			s:    NewSelector[TestModel](db).From("`test_model_t`").Where(C("Id").EQ(1)),
			want: &Query{
				SQL:  "SELECT * FROM `test_model_t` WHERE `Id` = ?;",
				Args: []any{1},
			},
		},
		{
			// 多个 predicate
			name: "multiple predicates",
			s: NewSelector[TestModel](db).
				Where(C("Age").GT(18), C("Age").LT(35)),
			want: &Query{
				SQL:  "SELECT * FROM `TestModel` WHERE (`Age` > ?) AND (`Age` < ?);",
				Args: []any{18, 35},
			},
		},
		{
			// 使用 AND
			name: "and",
			s: NewSelector[TestModel](db).
				Where(C("Age").GT(18).And(C("Age").LT(35))),
			want: &Query{
				SQL:  "SELECT * FROM `TestModel` WHERE (`Age` > ?) AND (`Age` < ?);",
				Args: []any{18, 35},
			},
		},
		{
			// 使用 OR
			name: "or",
			s:    NewSelector[TestModel](db).Where(C("Age").GT(18).Or(C("Age").LT(35))),
			want: &Query{
				SQL:  "SELECT * FROM `TestModel` WHERE (`Age` > ?) OR (`Age` < ?);",
				Args: []any{18, 35},
			},
		},
		{
			// 使用 NOT
			name: "not",
			s:    NewSelector[TestModel](db).Where(Not(C("Age").GT(18))),
			want: &Query{
				// NOT 前面有两个空格，因为我们没有对 NOT 进行特殊处理
				SQL:  "SELECT * FROM `TestModel` WHERE  NOT (`Age` > ?);",
				Args: []any{18},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.s.Build()
			assert.Equal(t, tt.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, tt.want, got)
		})
	}
}

type TestModel struct {
	Id        int64
	FirstName string
	Age       int8
	LastName  *sql.NullString
}

func TestSelector_Get(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	testCases := []struct {
		name     string
		query    string
		mockErr  error
		mockRows *sqlmock.Rows
		wantErr  error
		wantVal  *TestModel
	}{
		{
			name:    "single row",
			query:   "SELECT .*",
			mockErr: nil,
			mockRows: func() *sqlmock.Rows {
				rows := sqlmock.NewRows([]string{"id", "first_name", "age", "last_name"})
				rows.AddRow([]byte("123"), []byte("Ming"), []byte("18"), []byte("Deng"))
				return rows
			}(),
			wantVal: &TestModel{
				Id:        123,
				FirstName: "Ming",
				Age:       18,
				LastName:  &sql.NullString{Valid: true, String: "Deng"},
			},
		},

		{
			// SELECT 出来的行数小于你结构体的行数
			name:    "less columns",
			query:   "SELECT .*",
			mockErr: nil,
			mockRows: func() *sqlmock.Rows {
				rows := sqlmock.NewRows([]string{"id", "first_name"})
				rows.AddRow([]byte("123"), []byte("Ming"))
				return rows
			}(),
			wantVal: &TestModel{
				Id:        123,
				FirstName: "Ming",
			},
		},

		{
			name:    "invalid columns",
			query:   "SELECT .*",
			mockErr: nil,
			mockRows: func() *sqlmock.Rows {
				rows := sqlmock.NewRows([]string{"id", "first_name", "gender"})
				rows.AddRow([]byte("123"), []byte("Ming"), []byte("male"))
				return rows
			}(),
			wantErr: errs.NewErrUnknownColumn("gender"),
		},

		{
			name:    "more columns",
			query:   "SELECT .*",
			mockErr: nil,
			mockRows: func() *sqlmock.Rows {
				rows := sqlmock.NewRows([]string{"id", "first_name", "age", "last_name", "first_name"})
				rows.AddRow([]byte("123"), []byte("Ming"), []byte("18"), []byte("Deng"), []byte("明明"))
				return rows
			}(),
			wantErr: errs.ErrTooManyReturnColumns,
		},
	}
	for _, tc := range testCases {
		if tc.mockErr != nil {
			mock.ExpectQuery(tc.query).WillReturnError(tc.mockErr)
		} else {
			mock.ExpectQuery(tc.query).WillReturnRows(tc.mockRows)
		}
	}
	db, err := OpenDB(mockDB)
	require.NoError(t, err)
	for _, tt := range testCases {
		res, err := NewSelector[TestModel](db).Get(context.Background())
		assert.Equal(t, tt.wantErr, err)
		if err != nil {
			return
		}
		assert.Equal(t, tt.wantVal, res)
	}
}
