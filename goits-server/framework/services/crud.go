//go:generate mockgen -source=crud.go -destination=service_mock.go -package=mocking
package services

import (
	"context"

	"github.com/cpekyaman/goits/framework/orm"
	"github.com/cpekyaman/goits/framework/validation"
)

type GetAllService interface {
	GetAll(ctx context.Context) (interface{}, error)
	GetAllPaged(ctx context.Context, limit uint, offset uint64) (interface{}, error)
}
type GetByIdService interface {
	GetById(ctx context.Context, id uint64) (interface{}, error)
}
type ReaderService interface {
	GetAllService
	GetByIdService
}

type CreateService interface {
	Create(ctx context.Context, binding ObjectBinder) error
}
type UpdateService interface {
	Update(ctx context.Context, id uint64, binding ObjectBinder) error
}
type WriterService interface {
	CreateService
	UpdateService
}

func NewCrudService(repo orm.Repository) CrudServiceImpl {
	return CrudServiceImpl{repo}
}

type CrudService interface {
	ReaderService
	WriterService
}

type CrudServiceImpl struct {
	crudRepo orm.Repository
}

func (this CrudServiceImpl) Create(ctx context.Context, binding ObjectBinder, fullTypeName string, target orm.Entity) error {
	err := binding.BindTo(target)
	if err != nil {
		return err
	}

	if err := validation.ValidateStruct(fullTypeName, target); err != nil {
		return err
	}

	return this.crudRepo.Save(ctx, target)
}

func (this CrudServiceImpl) Update(ctx context.Context, id uint64, binding ObjectBinder, fullTypeName string, target orm.Entity) error {
	err := binding.BindTo(target)
	if err != nil {
		return err
	}

	target.SetId(id)

	if err := validation.ValidateStruct(fullTypeName, target); err != nil {
		return err
	}

	return this.crudRepo.Save(ctx, target)
}

type ObjectBinder interface {
	BindTo(target interface{}) error
}

// ObjectBinderFunc is used to bind input values to target object.
type ObjectBinderFunc func(target interface{}) error

func (obf ObjectBinderFunc) BindTo(target interface{}) error {
	return obf(target)
}
