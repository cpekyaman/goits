package metadata

import (
	"reflect"
	"sort"
)

var columnMapperRegistry map[string]ColumnMapper

func init() {
	columnMapperRegistry = make(map[string]ColumnMapper)
}

// ColumnMapper is responsible for mapping a field name to a column name.
type ColumnMapper interface {
	HasColumn(field string) bool

	Column(field string) string

	Fields() []string

	Columns() []string
}

// fieldMapColumnMapper is a ColumnMapper that uses a simple field-column map.
type fieldMapColumnMapper struct {
	fieldMap       map[string]string
	orderedFields  []string
	orderedColumns []string
}

func (this fieldMapColumnMapper) HasColumn(field string) bool {
	_, ok := this.fieldMap[field]
	return ok
}

func (this fieldMapColumnMapper) Column(field string) string {
	return this.fieldMap[field]
}

func (this fieldMapColumnMapper) Fields() []string {
	return this.orderedFields
}

func (this fieldMapColumnMapper) Columns() []string {
	return this.orderedColumns
}

// GetColumnMapper returns an already registered ColumnMapper for the entity represented by provided metadata.
func GetColumnMapper(ed EntityDef) (ColumnMapper, bool) {
	cm, found := columnMapperRegistry[ed.Name()]
	return cm, found
}

// NewColumnMapper creates a new ColumnMapper by introspecting the given example entity.
func NewColumnMapper(ed EntityDef, introspect interface{}) ColumnMapper {
	v := reflect.ValueOf(introspect)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	fieldMap := make(map[string]string)
	if v.Kind() == reflect.Struct {
		buildFieldMap(fieldMap, v)
	}

	// trying to have a deterministic order of columns in generated statements
	orderedFields := make([]string, len(fieldMap))
	orderedColumns := make([]string, len(fieldMap))
	i := 0
	for k, v := range fieldMap {
		orderedFields[i] = k
		orderedColumns[i] = v
		i++
	}
	sort.Strings(orderedFields)
	sort.Strings(orderedColumns)

	cm := fieldMapColumnMapper{fieldMap, orderedFields, orderedColumns}
	columnMapperRegistry[ed.Name()] = cm
	return cm
}

// buildFieldMap builds a map of field names and their corresponding column names by reflecting on value.
func buildFieldMap(fieldMap map[string]string, v reflect.Value) {
	for i := 0; i < v.NumField(); i++ {
		f := v.Type().Field(i)

		switch f.Type.Kind() {
		case reflect.Struct:
			if f.Anonymous {
				buildFieldMap(fieldMap, v.Field(i))
			} else {
				addField(f, fieldMap)
			}
		case reflect.Ptr:
			buildFieldMap(fieldMap, v.Elem())
		default:
			addField(f, fieldMap)
		}
	}
}

// addField adds the mapping of given field to fieldMap.
func addField(f reflect.StructField, fieldMap map[string]string) {
	fieldName := f.Name
	colName, found := f.Tag.Lookup("db")
	if found {
		fieldMap[fieldName] = colName
	} else {
		fieldMap[fieldName] = fieldName
	}
}
