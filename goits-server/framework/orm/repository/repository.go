package repository

import (
	"context"

	"github.com/cpekyaman/goits/framework/orm/domain"
	"github.com/cpekyaman/goits/framework/orm/metadata"
	"github.com/cpekyaman/goits/framework/orm/query"
	"github.com/cpekyaman/goits/framework/orm/db"
)

// ReaderRepository provides common querying functionality for db entities.
type ReaderRepository interface {
	// FindOneById finds a single entity by the given id.
	FindOneById(ctx context.Context, dest interface{}, id uint64) error

	// FindAll finds all entities.
	FindAll(ctx context.Context, dest interface{}) error

	// FindAllPaged finds all entities in the given page offset and limit.
	FindAllPaged(ctx context.Context, dest interface{}, limit uint, offset uint64) error

	// FindOneByAttribute is a generic finder to find a single entity by a unique attribute.
	FindOneByAttribute(ctx context.Context, dest interface{}, attr string, bindval interface{}) error

	// FindAllByAttributes finds all entities that match the criteria.
	FindAllByAttributes(ctx context.Context, dest interface{}, attrs map[string]interface{}) error

	// FindAllByAttributesPaged finds all entities in the given page offset and limit which match the criteria.
	FindAllByAttributesPaged(ctx context.Context, dest interface{}, attrs map[string]interface{}, limit uint, offset uint64) error
}

// WriterRepository provides basic modify functionality for db entities.
type WriterRepository interface {
	Save(ctx context.Context, entity domain.Entity) error
	Delete(ctx context.Context, id uint64) error
}

// Repository is a generic repository definition for standard crud operations.
type Repository interface {
	ReaderRepository
	WriterRepository
}

// NewRepository creates a new crud repository for given entity.
func NewRepository(ed metadata.EntityDef, introspect interface{}) SqlRepository {
	cm := metadata.NewColumnMapper(ed, introspect)
	qd := query.BuildQueryDef(introspect, ed, cm)

	return SqlRepository{db.DB(), ed, qd, cm}
}
