//go:generate mockgen -source=repository.go -destination=repository_mock.go -package=project
package project

import (
	"github.com/cpekyaman/goits/framework/orm"
)

var projectED orm.EntityDef

func init() {
	projectED = orm.NewEntityDef("Project", "project", "Id", "name asc", false)
}

type ProjectRepository interface {
	orm.Repository
}

type projectSqlRepository struct {
	orm.SqlRepository
}

func newProjectRepository() ProjectRepository {
	return projectSqlRepository{orm.NewRepository(&projectED, &Project{})}
}
