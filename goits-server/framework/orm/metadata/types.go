package metadata

import (
	"fmt"

	"github.com/cpekyaman/goits/config"
)

// ReadMetadataConfig reads the orm config file for the given module into memory.
func ReadMetadataConfig(module string) (map[string]ormEntityDef, error) {
	ed := make(map[string]ormEntityDef)
	err := config.ReadConfig(module, "orm", fmt.Sprintf("%s.orm", module), &ed)
	return ed, err
}

// EntityDef is the metadata for an entity to help customizing the repository behaviour in a generic way.
type EntityDef interface {
	Name() string
	Schema() string
	Table() string
	FullTableName() string
	PKColumn() string
	DefaultSort() string
	SoftDelete() bool
}

// ormEntityDef is the package private implementation for EntityDef.
type ormEntityDef struct {
	Name_        string `mapstructure:"name"`
	Schema_      string `mapstructure:"schema"`
	Table_       string `mapstructure:"table"`
	PkColumn_    string `mapstructure:"pkColumn"`
	DefaultSort_ string `mapstructure:"defaultSort"`
	SoftDelete_  bool   `mapstructure:"softDelete"`
}

func (this ormEntityDef) Name() string {
	return this.Name_
}
func (this ormEntityDef) Schema() string {
	return this.Schema_
}
func (this ormEntityDef) Table() string {
	return this.Table_
}
func (this ormEntityDef) FullTableName() string {
	return fmt.Sprintf("%s.%s", this.Schema_, this.Table_)
}
func (this ormEntityDef) PKColumn() string {
	return this.PkColumn_
}
func (this ormEntityDef) DefaultSort() string {
	return this.DefaultSort_
}
func (this ormEntityDef) SoftDelete() bool {
	return this.SoftDelete_
}
