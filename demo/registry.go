package demo

import (
	"reflect"
	"sync"
	"unicode"
	"zorm/demo/internal/errs"
)

type registry struct {
	models sync.Map
}

// get 这种方案
func (r *registry) get(val any) (*model, error) {
	typ := reflect.TypeOf(val)
	m, ok := r.models.Load(typ)
	if !ok {
		var err error
		if m, err = r.parseModel(typ); err != nil {
			return nil, err
		}
	}
	r.models.Store(typ, m)
	return m.(*model), nil
}

// parseModel 输入不能为 nil
func (r *registry) parseModel(typ reflect.Type) (*model, error) {
	if typ.Kind() != reflect.Ptr || typ.Elem().Kind() != reflect.Struct {
		return nil, errs.ErrPointerOnly
	}
	typ = typ.Elem()
	numField := typ.NumField()
	fieldMap := make(map[string]*field, numField)
	for i := 0; i < numField; i++ {
		fd := typ.Field(i)
		fieldMap[fd.Name] = &field{
			colName: underscoreName(fd.Name),
		}
	}
	return &model{
		tableName: underscoreName(typ.Name()),
		fieldMap:  fieldMap,
	}, nil
}

// underscoreName 驼峰转字符串命名
func underscoreName(tableName string) string {
	var buf []byte
	for i, v := range tableName {
		if unicode.IsUpper(v) {
			if i != 0 {
				buf = append(buf, '_')
			}
			buf = append(buf, byte(unicode.ToLower(v)))
		} else {
			buf = append(buf, byte(v))
		}

	}
	return string(buf)
}
