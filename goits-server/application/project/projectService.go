package project

import (
	"context"

	"github.com/cpekyaman/goits/framework/caching"
	"github.com/cpekyaman/goits/framework/services"
	"github.com/cpekyaman/goits/framework/validation"
)

type ProjectService interface {
	services.CRUDService
}

type projectServiceImpl struct {
	repo    ProjectRepository
	svcImpl services.CRUDServiceImpl
}

func newDefaultProjectService() ProjectService {
	return newProjectService(newProjectRepository(), caching.NamedCache("project"), validation.Provider())
}

func newProjectService(pr ProjectRepository, c caching.Cache, vp validation.ValidationProvider) ProjectService {
	return projectServiceImpl{pr, services.NewCRUDService(pr, c, vp)}
}

func (this projectServiceImpl) GetAll(ctx context.Context) (interface{}, error) {
	var resultList []Project
	err := this.repo.FindAll(ctx, &resultList)
	return resultList, err
}

func (this projectServiceImpl) GetAllPaged(ctx context.Context, limit uint, offset uint64) (interface{}, error) {
	var resultList []Project
	err := this.repo.FindAllPaged(ctx, &resultList, limit, offset)
	return resultList, err
}

func (this projectServiceImpl) GetById(ctx context.Context, id uint64) (interface{}, error) {
	result, err := this.svcImpl.Cache().GetOrCompute(caching.IdToKey(id), func() (interface{}, error) {
		var result Project
		err := this.repo.FindOneById(ctx, &result, id)
		return &result, err
	})

	return result, err
}

func (this projectServiceImpl) FindOne(ctx context.Context, attr string, attrValue interface{}) (interface{}, error) {
	var result Project
	err := this.repo.FindOneByAttribute(ctx, &result, attr, attrValue)
	return result, err
}

func (this projectServiceImpl) FindAll(ctx context.Context, attrs map[string]interface{}) (interface{}, error) {
	var resultList []Project
	err := this.repo.FindAllByAttributes(ctx, &resultList, attrs)
	return resultList, err
}

func (this projectServiceImpl) FindAllPaged(ctx context.Context, attrs map[string]interface{}, limit uint, offset uint64) (interface{}, error) {
	var resultList []Project
	err := this.repo.FindAllByAttributesPaged(ctx, &resultList, attrs, limit, offset)
	return resultList, err
}

func (this projectServiceImpl) Create(ctx context.Context, binding services.ObjectBinder) error {
	return this.svcImpl.Create(ctx, binding, projectTypeName, &Project{})
}

func (this projectServiceImpl) Update(ctx context.Context, id uint64, binding services.ObjectBinder) error {
	return this.svcImpl.Update(ctx, id, binding, projectTypeName, &Project{})
}

func (this projectServiceImpl) Delete(ctx context.Context, id uint64) error {
	return this.svcImpl.Delete(ctx, id)
}
