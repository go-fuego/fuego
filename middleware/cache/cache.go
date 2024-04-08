package cache

import (
	"bytes"
	"net/http"
	"time"
)

type Storage interface {
	Get(key string) (string, bool)
	Set(key, value string)
}

type Config struct {
	Storage Storage
	Key     func(r *http.Request) string // Key returns the cache key for the request
}

// Cache the response of GET requests for a given duration.
// By default, it will use an in-memory cache with a duration of 3 seconds, with a maximum of 1000 entries.
// You can provide your own storage implementation by passing a Config struct to the middleware.
// You can also provide your own key function to generate the cache key for a given request.
// Headers can be used to invalidate the cache:
//   - Cache-Control: no-cache will bypass the cache
//   - Cache-Control: no-store might use the cache but will not store the response in the cache
func New(config ...Config) func(http.Handler) http.Handler {
	if len(config) > 1 {
		panic("Only one config is allowed")
	}

	c := Config{
		Storage: NewInMemoryCache(3*time.Second, 1000),
		Key: func(r *http.Request) string {
			return "httpcache_" + r.URL.Path + "_" + r.Header.Get("Content-Type")
		},
	}

	if len(config) == 1 {
		if config[0].Storage != nil {
			c.Storage = config[0].Storage
		}

		if config[0].Key != nil {
			c.Key = config[0].Key
		}
	}

	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodGet {
				h.ServeHTTP(w, r)
				return
			}

			if r.Header.Get("Cache-Control") == "no-cache" {
				h.ServeHTTP(w, r)
				return
			}

			key := c.Key(r)

			if val, ok := c.Storage.Get(key); ok {
				respContentType, ok := c.Storage.Get(key + "_response-content-type")
				if ok {
					w.Header().Set("Content-Type", respContentType)
				}

				w.WriteHeader(http.StatusOK)
				w.Header().Set("Cache", "hit")
				_, _ = w.Write([]byte(val))
				return
			}

			if r.Header.Get("Cache-Control") == "no-store" {
				h.ServeHTTP(w, r)
				return
			}

			buf := &bytes.Buffer{}
			multiWriter := &MultiHTTPWriter{
				ResponseWriter: w,
				cacheWriter:    buf,
			}

			h.ServeHTTP(multiWriter, r)

			c.Storage.Set(key, multiWriter.cacheWriter.(*bytes.Buffer).String())
			c.Storage.Set(key+"_response-content-type", multiWriter.Header().Get("Content-Type"))
		})
	}
}
