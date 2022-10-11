package demo

import (
	"database/sql"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSelector_Build(t *testing.T) {
	tests := []struct {
		name    string
		s       QueryBuilder
		want    *Query
		wantErr error
	}{
		{
			name: "from",
			s:    NewSelector[TestModel]().From("test_model_tab"),
			want: &Query{
				SQL: "SELECT * FROM test_model_tab;",
			},
		},
		{
			name: "no from",
			s:    NewSelector[TestModel](),
			want: &Query{
				SQL: "SELECT * FROM `TestModel`;",
			},
		},
		{
			name: "empty from",
			s:    NewSelector[TestModel]().From(""),
			want: &Query{
				SQL: "SELECT * FROM `TestModel`;",
			},
		}, {
			name: "sigle and simple predicate",
			s:    NewSelector[TestModel]().From("`test_model_t`").Where(C("Id").EQ(1)),
			want: &Query{
				SQL:  "SELECT * FROM `test_model_t` WHERE `Id` = ?;",
				Args: []any{1},
			},
		},
		{
			// 多个 predicate
			name: "multiple predicates",
			s: NewSelector[TestModel]().
				Where(C("Age").GT(18), C("Age").LT(35)),
			want: &Query{
				SQL:  "SELECT * FROM `TestModel` WHERE (`Age` > ?) AND (`Age` < ?);",
				Args: []any{18, 35},
			},
		},
		{
			// 使用 AND
			name: "and",
			s: NewSelector[TestModel]().
				Where(C("Age").GT(18).And(C("Age").LT(35))),
			want: &Query{
				SQL:  "SELECT * FROM `TestModel` WHERE (`Age` > ?) AND (`Age` < ?);",
				Args: []any{18, 35},
			},
		},
		{
			// 使用 OR
			name: "or",
			s: NewSelector[TestModel]().
				Where(C("Age").GT(18).Or(C("Age").LT(35))),
			want: &Query{
				SQL:  "SELECT * FROM `TestModel` WHERE (`Age` > ?) OR (`Age` < ?);",
				Args: []any{18, 35},
			},
		},
		{
			// 使用 NOT
			name: "not",
			s:    NewSelector[TestModel]().Where(Not(C("Age").GT(18))),
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
