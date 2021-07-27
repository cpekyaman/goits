package project

import (
	"github.com/cpekyaman/goits/framework/orm/domain"
	"github.com/cpekyaman/goits/framework/validation"
)

const (
	projectTypeName       = "project.Project"
	projectStatusTypeName = "project.ProjectStatus"
	projectTypeTypeName   = "project.ProjectType"
)

func initDomain() {
	registerProjectValidations()
}

type ProjectType struct {
	domain.VersionedEntity
	Name        string `json:"name" db:"name"`
	Description string `json:"desc" db:"description"`
}

type ProjectStatus struct {
	domain.VersionedEntity
	Name        string `json:"name" db:"name"`
	Description string `json:"desc" db:"description"`
}

type Project struct {
	domain.VersionedTimeStampedEntity
	Name        string `json:"name" db:"name"`
	Description string `json:"desc" db:"description"`
	Type        uint64 `json:"type" db:"type"`
	Status      uint64 `json:"status" db:"status"`
}

func registerProjectValidations() {
	sv := validation.Struct(projectTypeName)

	sv.WithMandatoryName().
		WithMandatoryDesc().
		Field("type").With(validation.ValidId()).
		Field("status").With(validation.ValidId())
}

func NewProject() *Project {
	return &Project{domain.VersionedTimeStampedEntity{}, "", "", 0, 0}
}
