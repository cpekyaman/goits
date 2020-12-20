package caching

// ShardedCache is a type of Cache that uses sharding to improve performance.
// A ShardedCache is essentially a container of CacheShards which actually keep the values.
type ShardedCache interface {
	Cache
	getShards() []CacheShard
}

// CacheShard represents a shard of a ShardedCache, which actuall keeps part of the data.
type CacheShard interface {
	// Add adds a new value to the shard for the given key.
	Add(key string, item interface{}) bool

	// Get retrieves the value associated with the key, if it exists.
	Get(key string) (interface{}, bool)

	// Remove removes the value associated with the given key from cache.
	Remove(key string)

	// RemoveAll removes all values from the shard.
	RemoveAll()
}