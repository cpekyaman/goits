//go:generate mockgen -source=cache.go -destination=cache_mock.go -package=mocking
package caching

import (
	"hash/fnv"
	"time"
)

// CacheConfig provides a simple means to configure a cache upon creation.
type CacheConfig struct {
	Name        string
	MaxElements uint32
	TTLSeconds  uint32
}

// CacheProvider is responsible for creating a new instance of a specific type of cache.
type CacheProvider interface {
	NewCache(config CacheConfig) Cache
}

var cacheProvider CacheProvider

func init() {
	cacheProvider = MemoryCacheProvier{}
}

// CacheEntry is a simple wrapper around the actual value to be put in the cache.
type CacheEntry struct {
	Value     interface{}
	Timestamp time.Time
}

// Cache is a simple abstraction of a caching system.
type Cache interface {
	// Put adds a new value to the cache or replaces the existing value for the key.
	Put(key string, value interface{}) bool

	// Get retrieves the value associated with the key, if it exists.
	Get(key string) (interface{}, bool)

	// GetOrCompute tries to get the value if exists, similar to Get.
	// If value is not present, it runs the compute function and puts the result in the cache and returns.
	GetOrCompute(key string, compute func() (interface{}, error)) (interface{}, error)

	// InvalidateAll removes all items from the cache.
	InvalidateAll()

	// Invalidate removes the value associated with the key from the cache.
	Invalidate(key string)
}

// Provider returns the default CacheProvider the application should use.
func Provider() CacheProvider {
	return cacheProvider
}

// shardFor finds the target shard for the given key.
func shardFor(sc ShardedCache, key string) CacheShard {
	hash := fnv.New32()
	hash.Write([]byte(key))
	idx := hash.Sum32()
	return sc.getShards()[idx%uint32(shardCount)]
}
