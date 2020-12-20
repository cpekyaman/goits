package caching

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

type TestData struct {
	id   uint64
	name string
}

var cp CacheProvider
var cc CacheConfig

func init() {
	cp = MemoryCacheProvier{}
	cc = CacheConfig{"test", 100, 600}
}

func TestPutAndGet(t *testing.T) {
	cache := cp.NewCache(cc)

	testData := TestData{
		id:   1,
		name: "cenk",
	}

	b := cache.Put(testData.name, &testData)
	assert.True(t, b, "could not put item into cache")

	item, b := cache.Get(testData.name)
	assert.True(t, b, "could not get item from cache")

	itemPtr := item.(*TestData)
	itemPtr.name = "cenk updated"
	assert.Equal(t, "cenk updated", testData.name, "did not get back original item")
}

func TestGetOrCompute(t *testing.T) {
	cache := cp.NewCache(cc)

	var key = "100"
	var value = TestData{
		id:   1,
		name: "cenk",
	}

	item, err := cache.GetOrCompute(key, func() (interface{}, error) {
		return &value, nil
	})
	assert.Nil(t, err, "no error should be returned")

	itemPtr := item.(*TestData)
	itemPtr.name = "cenk changed"
	assert.Equal(t, "cenk changed", value.name, "original value should be updated")

	_, found := cache.Get(key)
	assert.True(t, found, "we should find the item second time")

	item, err = cache.GetOrCompute(key, func() (interface{}, error) {
		return nil, fmt.Errorf("not to be executed")
	})
	assert.Nil(t, err, "compute should not be executed second time")
}

func TestInvaliateAll(t *testing.T) {
	cache := cp.NewCache(cc)

	for i := 0; i < 10; i++ {
		cache.Put(strconv.Itoa(i), &TestData{uint64(i), "test"})
	}
	for i := 0; i < 10; i++ {
		_, found := cache.Get(strconv.Itoa(i))
		assert.True(t, found, "could not found item")
	}

	cache.InvalidateAll()
	for i := 0; i < 10; i++ {
		_, found := cache.Get(strconv.Itoa(i))
		assert.False(t, found, "should not find any item after invalidate")
	}
	for i := 0; i < 10; i++ {
		cache.Put(strconv.Itoa(i), &TestData{uint64(i), "test"})
		_, found := cache.Get(strconv.Itoa(i))
		assert.True(t, found, "cache did not work as expected after invalidate")
	}
}

func TestInvalidate(t *testing.T) {
	cache := cp.NewCache(cc)

	var key = "100"
	var item = TestData{1, "test data"}

	cache.Put(key, &item)

	_, found := cache.Get(key)
	assert.True(t, found, "should've found item before invalidate")

	cache.Invalidate(key)
	_, found = cache.Get(key)
	assert.False(t, found, "should not find item after invalidate")
}
