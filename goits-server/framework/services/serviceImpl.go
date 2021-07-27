package services

import (
	"context"

	"github.com/cpekyaman/goits/framework/caching"
	"github.com/cpekyaman/goits/framework/orm/repository"
	"github.com/cpekyaman/goits/framework/orm/domain"
	"github.com/cpekyaman/goits/framework/validation"
)

// CRUDServiceImpl is a helper implementation class for crud services.
type CRUDServiceImpl struct {
	crudRepo repository.Repository
	cache    caching.Cache
	vp       validation.ValidationProvider
}

func (this CRUDServiceImpl) Cache() caching.Cache {
	return this.cache
}

// Create binds input data to target entity by using provided binding, performs validations and saves the new entity.
func (this CRUDServiceImpl) Create(ctx context.Context, binding ObjectBinder, fullTypeName string, target domain.Entity) error {
	err := binding.BindTo(target)
	if err != nil {
		return err
	}

	if err := this.vp.ValidateStruct(fullTypeName, target); err != nil {
		return err
	}

	return this.crudRepo.Save(ctx, target)
}

// Update binds input data to target entity by using provided binding, performs validations and saves the updated entity.
func (this CRUDServiceImpl) Update(ctx context.Context, id uint64, binding ObjectBinder, fullTypeName string, target domain.Entity) error {
	err := this.crudRepo.FindOneById(ctx, target, id)
	if err != nil {
		return err
	}

	err = binding.BindTo(target)
	if err != nil {
		return err
	}

	if err := this.vp.ValidateStruct(fullTypeName, target); err != nil {
		return err
	}

	err = this.crudRepo.Save(ctx, target)
	if err == nil && this.cache != nil {
		this.cache.Invalidate(caching.IdToKey(id))
	}
	return err
}

// Delete simply deletes the entity represented by the given id.
func (this CRUDServiceImpl) Delete(ctx context.Context, id uint64) error {
	err := this.crudRepo.Delete(ctx, id)
	if err == nil && this.cache != nil {
		this.cache.Invalidate(caching.IdToKey(id))
	}
	return err
}