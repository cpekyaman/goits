package caching

import (
	"sync"
	"time"
)

const (
	shardCount = 16
)

// ShardedMemoryCacheProvier is a CacheProvider that creates in memory sharded caches.
type ShardedMemoryCacheProvier struct{}

// NewCache creates a new ShardedMemoryCache and configures it by using given configuration.
func (cp ShardedMemoryCacheProvier) NewCache(config CacheConfig) Cache {
	sa := make([]CacheShard, shardCount)
	for i := 0; i < shardCount; i++ {
		sa[i] = &memoryCacheShard{items: make(map[string]CacheEntry)}
	}
	return &ShardedMemoryCache{config: config, shards: sa}
}

// ShardedMemoryCache is a ShardedCache implementation that keeps data in local memory
type ShardedMemoryCache struct {
	config CacheConfig
	shards []CacheShard
}

func (this *ShardedMemoryCache) Put(key string, value interface{}) bool {
	s := shardFor(this, key)
	return s.Add(key, value)
}

func (this *ShardedMemoryCache) Get(key string) (interface{}, bool) {
	s := shardFor(this, key)
	return s.Get(key)
}

func (this *ShardedMemoryCache) GetOrCompute(key string, compute func() (interface{}, error)) (interface{}, error) {
	s := shardFor(this, key)
	result, found := s.Get(key)

	if !found {
		item, err := compute()
		if err == nil {
			s.Add(key, item)
			return item, nil
		}
		return nil, err
	}

	return result, nil
}

func (this *ShardedMemoryCache) InvalidateAll() {
	for _, s := range this.shards {
		s.RemoveAll()
	}
}

func (this *ShardedMemoryCache) Invalidate(key string) {
	s := shardFor(this, key)
	s.Remove(key)
}

func (this *ShardedMemoryCache) getShards() []CacheShard {
	return this.shards
}

type memoryCacheShard struct {
	items map[string]CacheEntry
	sync.RWMutex
}

func (s *memoryCacheShard) Add(key string, item interface{}) bool {
	s.Lock()
	s.items[key] = CacheEntry{item, time.Now()}
	s.Unlock()
	return true
}

func (s *memoryCacheShard) Get(key string) (interface{}, bool) {
	s.RLock()
	result, found := s.items[key]
	s.RUnlock()
	if !found {
		return nil, false
	}
	return result.Value, true
}

func (s *memoryCacheShard) Remove(key string) {
	s.Lock()
	delete(s.items, key)
	s.Unlock()
}

func (s *memoryCacheShard) RemoveAll() {
	s.Lock()
	s.items = make(map[string]CacheEntry)
	s.Unlock()
}
