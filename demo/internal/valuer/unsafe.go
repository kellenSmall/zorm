package valuer

import (
	"database/sql"
	"reflect"
	"unsafe"
	"zorm/demo/internal/errs"
	"zorm/demo/internal/model"
)

type valueUnsafe struct {
	t     any
	model *model.Model
	add   unsafe.Pointer
}

func (u valueUnsafe) SetColumns(rows *sql.Rows) error {
	if !rows.Next() {
		return errs.ErrNoRows
	}
	cols, err := rows.Columns()
	if err != nil {
		return err
	}

	if len(cols) > len(u.model.FieldMap) {
		return errs.ErrTooManyReturnColumns
	}
	colsValues := make([]any, 0, len(cols))
	for _, col := range cols {
		fd, ok := u.model.ColsMap[col]
		if !ok {
			return errs.NewErrUnknownColumn(col)
		}
		fdVal := reflect.NewAt(fd.Type, unsafe.Pointer(uintptr(u.add)+fd.Offset))
		colsValues = append(colsValues, fdVal.Interface())
	}
	return rows.Scan(colsValues...)
}

func NewValueUnsafe(t any, model *model.Model) Value {
	addr := unsafe.Pointer(reflect.ValueOf(t).Pointer())
	return valueUnsafe{
		t:     t,
		model: model,
		add:   addr,
	}
}
