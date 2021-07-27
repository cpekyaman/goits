package domain

import (
	"time"
)

// Entity is the base interface for entities to implement.
type Entity interface {
	GetId() uint64

	SetId(id uint64)
}

// DomainEntity is the default base type that entities compose from.
type DomainEntity struct {
	Id uint64 `json:"id" db:"id"`
}

func (this *DomainEntity) GetId() uint64 {
	return this.Id
}

func (this *DomainEntity) SetId(id uint64) {
	this.Id = id
}

// Versioned is the interface to implement if the entity uses opmitistic locking with version.
type Versioned interface {
	GetVersion() uint32
}

// VersionedEntity is the default base type that versioned entities compose from.
type VersionedEntity struct {
	Id      uint64 `json:"id" db:"id"`
	Version uint32 `json:"version" db:"version"`
}

func (this *VersionedEntity) GetId() uint64 {
	return this.Id
}

func (this *VersionedEntity) SetId(id uint64) {
	this.Id = id
}

func (this *VersionedEntity) GetVersion() uint32 {
	return this.Version
}

// Timestamped is the interface to implement if the entity keeps track of times it changes.
type Timestamped interface {
	GetCreateTime() time.Time
	GetLastModifiedTime() time.Time
}

// TimestampedEntity is the struct to embed for any entity that keeps track of create and update timestamps.
type TimestampedEntity struct {
	Id               uint64    `json:"id" db:"id"`
	CreateTime       time.Time `json:"createdAt" db:"create_time"`
	LastModifiedTime time.Time `json:"lastModifiedAt" db:"last_modified_time"`
}

func (this *TimestampedEntity) GetId() uint64 {
	return this.Id
}

func (this *TimestampedEntity) SetId(id uint64) {
	this.Id = id
}

func (this *TimestampedEntity) GetCreateTime() time.Time {
	return this.CreateTime
}

func (this *TimestampedEntity) GetLastModifiedTime() time.Time {
	return this.LastModifiedTime
}

// VersionedTimestamped is the interface to implement if both versioning and time tracking is supported.
type VersionedTimestamped interface {
	Versioned
	Timestamped
}

// VersionedTimeStampedEntity is the struct to embed for any versioned entity that keeps track of create and update timestamps.
type VersionedTimeStampedEntity struct {
	Id               uint64    `json:"id" db:"id"`
	Version          uint32    `json:"version" db:"version"`
	CreateTime       time.Time `json:"createdAt" db:"create_time"`
	LastModifiedTime time.Time `json:"lastModifiedAt" db:"last_modified_time"`
}

func (this *VersionedTimeStampedEntity) GetId() uint64 {
	return this.Id
}

func (this *VersionedTimeStampedEntity) SetId(id uint64) {
	this.Id = id
}

func (this *VersionedTimeStampedEntity) GetVersion() uint32 {
	return this.Version
}

func (this *VersionedTimeStampedEntity) GetCreateTime() time.Time {
	return this.CreateTime
}

func (this *VersionedTimeStampedEntity) GetLastModifiedTime() time.Time {
	return this.LastModifiedTime
}
