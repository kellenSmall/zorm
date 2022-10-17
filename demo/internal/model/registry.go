package model

import (
	"reflect"
	"strings"
	"sync"
	"unicode"
	"zorm/demo/internal/errs"
)

type Registry interface {
	Get(val any) (*Model, error)
	Register(val any, opts ...Option) (*Model, error)
}

type registry struct {
	models sync.Map
}

func NewRegistry() Registry {
	return &registry{}
}

// get 这种方案
func (r *registry) Get(val any) (*Model, error) {
	typ := reflect.TypeOf(val)
	m, ok := r.models.Load(typ)
	if ok {
		return m.(*Model), nil
	}

	return r.Register(val)
}

// parseModel 输入不能为 nil
func (r *registry) Register(val any, opts ...Option) (*Model, error) {
	if val == nil {
		return nil, errs.ErrInputNil
	}
	typ := reflect.TypeOf(val)
	if typ.Kind() != reflect.Ptr || typ.Elem().Kind() != reflect.Struct {
		return nil, errs.ErrPointerOnly
	}
	typ = typ.Elem()
	numField := typ.NumField()
	fieldMap := make(map[string]*Field, numField)
	colsMap := make(map[string]*Field, numField)
	for i := 0; i < numField; i++ {
		fd := typ.Field(i)
		tags, err := r.parseTag(fd.Tag)
		if err != nil {
			return nil, err
		}
		colName := tags["column"]
		if colName == "" {
			colName = underscoreName(fd.Name)
		}
		field := &Field{
			ColName: colName,
			GoName:  fd.Name,
			Type:    fd.Type,
			Offset:  fd.Offset,
		}
		fieldMap[fd.Name] = field
		colsMap[colName] = field
	}
	var tableName string
	if tn, ok := val.(TableName); ok {
		tableName = tn.TableName()
	}
	if tableName == "" {
		tableName = underscoreName(typ.Name())
	}
	res := &Model{
		TableName: tableName,
		FieldMap:  fieldMap,
		ColsMap:   colsMap,
	}
	for _, opt := range opts {
		if err := opt(res); err != nil {
			return nil, err
		}

	}
	r.models.Store(typ, res)
	return res, nil
}

func (r *registry) parseTag(tag reflect.StructTag) (map[string]string, error) {
	ormTag := tag.Get("orm")
	if ormTag == "" {
		return map[string]string{}, nil
	}
	res := make(map[string]string, 1)
	kvs := strings.Split(ormTag, ",")
	for _, pair := range kvs {
		kv := strings.Split(pair, "=")
		if len(kv) != 2 {
			return nil, errs.NewErrInvalidTagContent(pair)
		}
		res[kv[0]] = kv[1]
	}
	return res, nil
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
