package query

import (
	"fmt"
	"strings"

	"github.com/cpekyaman/goits/framework/orm/metadata"
)

// BuildQueryByAttributes builds a select query by using provided attribute map as criteria.
func BuildQueryByAttributes(ed metadata.EntityDef, qd QueryDef, cm metadata.ColumnMapper, attrs map[string]interface{}, limit uint, offset uint64) (string, []interface{}) {
	where, params := BuildCriteria(ed, qd, cm, attrs)

	template := findAllByAttributesTemplate
	if limit > 0 {
		template = findAllByAttributesPagedTemplate
	}

	return fmt.Sprintf(template, qd.SelectColumns(), ed.FullTableName(), where, ed.DefaultSort()), params
}

// BuildCriteria builds the where fragment of the select query by using provided attributes and their values.
func BuildCriteria(ed metadata.EntityDef, qd QueryDef, cm metadata.ColumnMapper, attrs map[string]interface{}) (string, []interface{}) {
	var criteria []string
	var params []interface{}
	var idx = 1
	for k, v := range attrs {
		if cm.HasColumn(k) {
			criteria = append(criteria, fmt.Sprintf("%s=$%d", cm.Column(k), idx))
			params = append(params, v)
			idx++
		}
	}
	return strings.Join(criteria, " AND "), params
}

// BuildFindOneQuery builds a single row select query by using attr as the only criteria.
// It is expected that attr corresponds to a unique column (by definition or in practice).
func BuildFindOneQuery(ed metadata.EntityDef, qd QueryDef, cm metadata.ColumnMapper, attr string) string {
	return fmt.Sprintf(findOneByAttributeTemplate, qd.SelectColumns(), ed.FullTableName(), cm.Column(attr))
}

// BuildFindAllPagedQuery builds the paging on top of default find all query for the given offset and limit values.
func BuildFindAllPagedQuery(ed metadata.EntityDef, qd QueryDef, limit uint, offset uint64) string {
	return fmt.Sprintf(findAllPagedTemplate, qd.SelectColumns(), ed.FullTableName(), ed.DefaultSort(), limit, offset)
}
