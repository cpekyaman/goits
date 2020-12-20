package project

import (
	"context"
	"strconv"

	"github.com/cpekyaman/goits/framework/caching"
	"github.com/cpekyaman/goits/framework/services"
)

type ProjectService interface {
	services.ReaderService
}

type defaultProjectService struct {
	cache   caching.Cache
	repo    ProjectRepository
	svcImpl services.CrudServiceImpl
}

func newProjectService() ProjectService {
	cc := caching.CacheConfig{
		Name:        "Project",
		MaxElements: 100,
		TTLSeconds:  3600,
	}

	return newCachingProjectService(caching.Provider().NewCache(cc))
}

func newCachingProjectService(c caching.Cache) ProjectService {
	pr := newProjectRepository()
	return defaultProjectService{c, pr, services.NewCrudService(pr)}
}

func (this defaultProjectService) GetAll(ctx context.Context) (interface{}, error) {
	var projects []Project
	err := this.repo.FindAll(ctx, &projects)
	return projects, err
}

func (this defaultProjectService) GetAllPaged(ctx context.Context, limit uint, offset uint64) (interface{}, error) {
	var projects []Project
	err := this.repo.FindAllPaged(ctx, &projects, limit, offset)
	return projects, err
}

func (this defaultProjectService) GetById(ctx context.Context, id uint64) (interface{}, error) {
	prj, err := this.cache.GetOrCompute(strconv.FormatUint(id, 10), func() (interface{}, error) {
		var prj Project
		err := this.repo.FindOneById(ctx, &prj, id)
		return &prj, err
	})

	return prj, err
}

func (this defaultProjectService) Create(ctx context.Context, binding services.ObjectBinder) error {
	return this.svcImpl.Create(ctx, binding, projectTypeName, &Project{})
}

func (this defaultProjectService) Update(ctx context.Context, id uint64, binding services.ObjectBinder) error {
	return this.svcImpl.Update(ctx, id, binding, projectTypeName, &Project{})
}
