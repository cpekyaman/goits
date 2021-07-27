//go:generate mockgen -source=repository.go -destination=repository_mock.go -package=project
package project

import (
	"github.com/cpekyaman/goits/framework/orm/domain"
	"github.com/cpekyaman/goits/framework/orm/metadata"
	"github.com/cpekyaman/goits/framework/orm/repository"
)

var projectStatusED metadata.EntityDef
var projectTypeED metadata.EntityDef
var projectED metadata.EntityDef

func init() {
	domain.RegisterEntityConfig("project")

	projectStatusED = domain.EntityDefByName(projectStatusTypeName)
	projectTypeED = domain.EntityDefByName(projectTypeTypeName)
	projectED = domain.EntityDefByName(projectTypeName)
}

type ProjectRepository interface {
	repository.Repository
}
type projectSqlRepository struct {
	repository.SqlRepository
}

func newProjectRepository() ProjectRepository {
	return projectSqlRepository{repository.NewRepository(projectED, &Project{})}
}

type ProjectTypeRepository interface {
	repository.Repository
}

type projectTypeSqlReporsitory struct {
	repository.SqlRepository
}

func newProjectTypeRepository() ProjectTypeRepository {
	return projectTypeSqlReporsitory{repository.NewRepository(projectTypeED, &ProjectType{})}
}

type ProjectStatusRepository interface {
	repository.Repository
}

type projectStatusSqlRepository struct {
	repository.SqlRepository
}

func newProjectStatusRepository() ProjectStatusRepository {
	return projectStatusSqlRepository{repository.NewRepository(projectStatusED, &ProjectStatus{})}
}
