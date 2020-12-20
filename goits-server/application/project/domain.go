package project

import (
	"github.com/cpekyaman/goits/framework/orm"
	"github.com/cpekyaman/goits/framework/validation"
)

const (
	projectTypeName       = "project.Project"
	projectStatusTypeName = "project.ProjectStatus"
)

func init() {
	registerProjectValidations()
}

type ProjectStatus struct {
	orm.VersionedEntity
	Name        string `json:"name" db:"name"`
	Description string `json:"desc" db:"description"`
}

type Project struct {
	orm.VersionedTimeStampedEntity
	Name        string `json:"name" db:"name"`
	Description string `json:"desc" db:"description"`
	Status      uint64 `json:"status" db:"status"`
}

func registerProjectValidations() {
	sv := validation.Struct(projectTypeName)

	sv.Field("name").
		With(validation.NotBlank(), validation.Pattern(validation.PatternAlNum)).
		Field("description").
		With(validation.NotBlank(), validation.Pattern(validation.PatternAlNum))
}

func applyProjectValidations(vc *validation.ValidationContext, p Project) error {
	vc.
		Validate("id", p.Id).
		Validate("name", p.Name).
		Validate("description", p.Description)

	return vc.Errors()
}

func NewProject() *Project {
	return &Project{orm.VersionedTimeStampedEntity{}, "", "", 0}
}
