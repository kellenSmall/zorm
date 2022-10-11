package demo

type model struct {
	tableName string
	fieldMap  map[string]*field
}
type field struct {
	colName string
}
