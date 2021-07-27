package metadata

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var ed EntityDef

func init() {
	ed = ormEntityDef{
		Name_: "IntrospectTestEntity",
	}
}

type IntrospectTestEntity struct {
	Name string `db:"name"`
	Age  uint32 `db:"age"`
}

func TestColumnMapper(t *testing.T) {
	// given
	target := &IntrospectTestEntity{}

	// when
	cm := NewColumnMapper(ed, target)

	// then
	assert.Equal(t, "Age", cm.Fields()[0], "first element of fields should be Age")
	assert.Equal(t, "age", cm.Columns()[0], "first element of columns should be age")

	colName := cm.Column("CreateTime")
	assert.Equal(t, "create_time", colName, "not correct column name")
}
