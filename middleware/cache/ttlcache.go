package cache

import (
	"time"

	"github.com/jellydator/ttlcache/v3"
)

type TTLCache struct {
	cache *ttlcache.Cache[string, string]
}

var _ Storage = (*TTLCache)(nil)

func NewTTL() *TTLCache {
	cache := ttlcache.New[string, string](
		ttlcache.WithTTL[string, string](30 * time.Minute),
	)

	go cache.Start() // starts automatic expired item deletion

	return &TTLCache{
		cache: cache,
	}
}

func (t *TTLCache) Get(key string) (string, bool) {
	val := t.cache.Get(key)

	if val != nil {
		return val.Value(), true
	}
	return "", false
}

func (t *TTLCache) Set(key string, value string) {
	t.cache.Set(key, value, 5*time.Second)
}

func (t *TTLCache) Delete(key string) {
	t.cache.Delete(key)
}
