//go:generate mockgen -source=service.go -destination=service_mock.go -package=mocking
package services

import (
	"context"

	"github.com/cpekyaman/goits/framework/caching"
	"github.com/cpekyaman/goits/framework/orm/repository"
	"github.com/cpekyaman/goits/framework/validation"
)

// SearcherService defines methods that support searching by attribute(s).
type SearcherService interface {
	FindOne(ctx context.Context, attr string, attrValue interface{}) (interface{}, error)
	FindAll(ctx context.Context, attrs map[string]interface{}) (interface{}, error)
	FindAllPaged(ctx context.Context, attrs map[string]interface{}, limit uint, offset uint64) (interface{}, error)
}

// ReaderService defines methods that are about fetching existing data.
type ReaderService interface {
	GetAll(ctx context.Context) (interface{}, error)
	GetAllPaged(ctx context.Context, limit uint, offset uint64) (interface{}, error)
	GetById(ctx context.Context, id uint64) (interface{}, error)
}

// ObjectBinder is used to bind input data (e.g. from an http request) to a target entity.
type ObjectBinder interface {
	BindTo(target interface{}) error
}

// ObjectBinderFunc is used to bind input values to target object.
type ObjectBinderFunc func(target interface{}) error

func (obf ObjectBinderFunc) BindTo(target interface{}) error {
	return obf(target)
}

// CreatorService defines the method to create a new entity.
type CreatorService interface {
	Create(ctx context.Context, binding ObjectBinder) error
}

// UpdaterService defines the method to update an existing entity.
type UpdaterService interface {
	Update(ctx context.Context, id uint64, binding ObjectBinder) error
}

// DeleterService defines the method to delete an existing entity.
type DeleterService interface {
	Delete(ctx context.Context, id uint64) error
}

// WriterService combines create and update methods into a single interface.
type WriterService interface {
	CreatorService
	UpdaterService
}

// ReaderWriterService combines read and write operations (except delete) into a single interface.
type ReaderWriterService interface {
	ReaderService
	WriterService
}

// CRUDService is the general interface that support all common crud operations.
type CRUDService interface {
	ReaderWriterService
	SearcherService
	DeleterService
}

// NewCRUDService creates a new CRUDServiceImpl that uses the provided repository for db operations.
func NewCRUDService(repo repository.Repository, cache caching.Cache, vp validation.ValidationProvider) CRUDServiceImpl {
	return CRUDServiceImpl{repo, cache, vp}
}
