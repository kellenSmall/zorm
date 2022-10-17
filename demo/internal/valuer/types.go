package valuer

import (
	"database/sql"
	"zorm/demo/internal/model"
)

type Value interface {
	SetColumns(rows *sql.Rows) error
}

type Creator func(t any, model *model.Model) Value
