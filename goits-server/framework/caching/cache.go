//go:generate mockgen -source=cache.go -destination=cache_mock.go -package=mocking
package caching

import (
	"hash/fnv"
	"strconv"
	"time"

	"github.com/cpekyaman/goits/config"
)

// CacheConfig provides a simple means to configure a cache upon creation.
type CacheConfig struct {
	Name        string `mapstructure:"name"`
	MaxElements uint32 `mapstructure:"maxElements"`
	TTLSeconds  uint32 `mapstructure:"ttlSeconds"`
}

// CacheProvider is responsible for creating a new instance of a specific type of cache.
type CacheProvider interface {
	NewCache(config CacheConfig) Cache
}

var cacheProvider CacheProvider
var cacheConfigs map[string]CacheConfig

func init() {
	cacheProvider = MemoryCacheProvier{}
	cacheConfigs = make(map[string]CacheConfig)
	config.ReadInto("caching", &cacheConfigs)
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

// NamedCache returns creates a cache by using pre defined named config.
// It returns a noop cache if named config is not found.
func NamedCache(name string) Cache {
	cc, ok := cacheConfigs[name]
	if !ok {
		return NoOpCache()
	} else {
		return cacheProvider.NewCache(cc)
	}
}

// CustomCache creates a cache by using the caller supplied cache config.
func CustomCache(config CacheConfig) Cache {
	return cacheProvider.NewCache(config)
}

// NoOpCache returns a noop cache which does no actually cache anything.
func NoOpCache() Cache {
	return noop
}

// IdToKey converts the id to string format to use it as a key.
func IdToKey(id uint64) string {
	return strconv.FormatUint(id, 10)
}

// shardFor finds the target shard for the given key.
func shardFor(sc ShardedCache, key string) CacheShard {
	hash := fnv.New32()
	hash.Write([]byte(key))
	idx := hash.Sum32()
	return sc.getShards()[idx%uint32(shardCount)]
}
