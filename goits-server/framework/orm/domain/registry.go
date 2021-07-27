package domain

import (
	"github.com/cpekyaman/goits/framework/monitoring"
	"github.com/cpekyaman/goits/framework/orm/metadata"
)

var entityMetaData map[string]metadata.EntityDef
var defaultNonInsertableColumns map[string]bool
var defaultNonUpdatableColumns map[string]bool

func init() {
	entityMetaData = make(map[string]metadata.EntityDef)

	defaultNonInsertableColumns = map[string]bool{
		"Id":               true,
		"Version":          true,
		"CreateTime":       true,
		"LastModifiedTime": true,
	}

	defaultNonUpdatableColumns = map[string]bool{
		"Id":               true,
		"Version":          true,
		"CreateTime":       true,
		"LastModifiedTime": true,
	}
}

// IsNonInsertableField checks if the given field is not insertable by default.
func IsNonInsertableField(f string) bool {
	return defaultNonInsertableColumns[f]
}

// IsNonUpdatableField checks if the given field is not updatable by default.
func IsNonUpdatableField(f string) bool {
	return defaultNonUpdatableColumns[f]
}

// EntityDefByName returns a registered entity metadata by given fully qualified type name.
func EntityDefByName(name string) metadata.EntityDef {
	return entityMetaData[name]
}

// RegisterEntityConfig reads the orm config file of the requested module and registers EntityDefs found there.
func RegisterEntityConfig(module string) {
	ed, err := metadata.ReadMetadataConfig(module)
	if err != nil {
		monitoring.RootLogger().With(monitoring.ErrLogField(err)).Fatal("could not find orm config for " + module)
	} else {
		for _, v := range ed {
			monitoring.RootLogger().With(monitoring.StrLogField("domainType", v.Name())).Info("Registering EntityDef")
			entityMetaData[v.Name()] = v
		}
	}
}
