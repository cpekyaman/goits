package caching

type noopCache struct{}

var noop noopCache

func (this noopCache) Put(key string, value interface{}) bool {
	return true
}

func (this noopCache) Get(key string) (interface{}, bool) {
	return nil, false
}

func (this noopCache) GetOrCompute(key string, compute func() (interface{}, error)) (interface{}, error) {
	r, err := compute()
	return r, err
}

func (this noopCache) InvalidateAll() {
	// nothing to invalidate
}

func (this noopCache) Invalidate(key string) {
	// nothing to invalidate
}
