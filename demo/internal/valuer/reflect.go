package valuer

import (
	"database/sql"
	"reflect"
	"zorm/demo/internal/errs"
	"zorm/demo/internal/model"
)

type valueReflect struct {
	t     any
	model *model.Model
}

func (v valueReflect) SetColumns(rows *sql.Rows) error {
	if !rows.Next() {
		return errs.ErrNoRows
	}
	cols, err := rows.Columns()
	if err != nil {
		return err
	}
	if len(cols) > len(v.model.FieldMap) {
		return errs.ErrTooManyReturnColumns
	}
	colValues := make([]any, 0, len(cols))
	colElemVals := make([]reflect.Value, 0, len(cols))
	for _, col := range cols {
		fd, ok := v.model.ColsMap[col]
		if !ok {
			return errs.NewErrUnknownColumn(col)
		}
		fdVal := reflect.New(fd.Type)
		colElemVals = append(colElemVals, fdVal.Elem())
		//rows.scan 需要 指针的变量
		colValues = append(colValues, fdVal.Interface())
	}
	err = rows.Scan(colValues...)
	if err != nil {
		return err
	}
	t := v.t
	//指针不能赋值 需要转成真是的结构体
	tVal := reflect.ValueOf(t).Elem()
	for i, col := range cols {
		fd := v.model.ColsMap[col]
		tVal.FieldByName(fd.GoName).Set(colElemVals[i])
	}
	return nil
}

var _ Value = valueReflect{}

func NewValueReflect(t any, model *model.Model) Value {
	return valueReflect{
		t:     t,
		model: model,
	}
}
