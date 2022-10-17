package model

import (
	"reflect"
	"zorm/demo/internal/errs"
)

type Option func(m *Model) error
type Model struct {
	TableName string
	FieldMap  map[string]*Field
	ColsMap   map[string]*Field
}
type Field struct {
	GoName  string
	ColName string
	Type    reflect.Type
	Offset  uintptr
}

const (
	tagKeyColumn = "column"
)

func WithTableName(name string) Option {
	return func(m *Model) error {
		// if name == "" {
		// 	return errs.ErrEmptyTableName
		// }
		m.TableName = name
		return nil
	}
}

func ModelWithColumnName(field string, colName string) Option {
	return func(m *Model) error {
		fd, ok := m.FieldMap[field]
		if !ok {
			return errs.NewErrUnknownColumn(field)
		}
		fd.ColName = colName
		return nil
	}
}

type TableName interface {
	TableName() string
}
