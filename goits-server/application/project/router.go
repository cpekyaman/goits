package project

import (
	"github.com/cpekyaman/goits/framework/routing"
)

type projectResource struct {
	svc ProjectService
	routing.ApiResource
}

var projectAPI projectResource

func InitProject() {
	projectAPI = newProjectResource(newProjectService())
	projectAPI.Register()
}

func newProjectResource(svc ProjectService) projectResource {
	res := projectResource{
		svc,
		routing.NewApiResource("Project", "project", svc),
	}

	return res
}
