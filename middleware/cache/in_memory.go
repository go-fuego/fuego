package cache

import (
	"time"

	"github.com/hashicorp/golang-lru/v2/expirable"
)

type TTLCache struct {
	cache *expirable.LRU[string, string]
}

var _ Storage = (*TTLCache)(nil)

func NewInMemoryCache(duration time.Duration, maxObjects int) *TTLCache {
	return &TTLCache{
		cache: expirable.NewLRU[string, string](maxObjects, nil, duration),
	}
}

func (t *TTLCache) Get(key string) (string, bool) {
	return t.cache.Get(key)
}

func (t *TTLCache) Set(key string, value string) {
	t.cache.Add(key, value)
}
