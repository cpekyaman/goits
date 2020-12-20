package caching

import (
	"sync"
	"time"
)

// MemoryCacheProvier is a CacheProvider that creates in memory caches.
type MemoryCacheProvier struct{}

// NewCache creates a new MemoryCache and configures it by using given configuration.
func (cp MemoryCacheProvier) NewCache(config CacheConfig) Cache {
	items := make(map[string]CacheEntry)
	return &MemoryCache{config: config, items: items}
}

// MemoryCache is a Cache implementation that keeps data in local memory
type MemoryCache struct {
	config CacheConfig
	items  map[string]CacheEntry
	sync.RWMutex
}

func (this *MemoryCache) Put(key string, value interface{}) bool {
	this.Lock()
	this.items[key] = CacheEntry{value, time.Now()}
	this.Unlock()
	return true
}

func (this *MemoryCache) Get(key string) (interface{}, bool) {
	this.RLock()
	result, found := this.items[key]
	this.RUnlock()
	if !found {
		return nil, false
	}
	return result.Value, true
}

func (this *MemoryCache) GetOrCompute(key string, compute func() (interface{}, error)) (interface{}, error) {
	result, found := this.Get(key)

	if !found {
		item, err := compute()
		if err == nil {
			this.Put(key, item)
			return item, nil
		}
		return nil, err
	}

	return result, nil
}

func (this *MemoryCache) InvalidateAll() {
	this.Lock()
	this.items = make(map[string]CacheEntry)
	this.Unlock()
}

func (this *MemoryCache) Invalidate(key string) {
	this.Lock()
	delete(this.items, key)
	this.Unlock()
}
