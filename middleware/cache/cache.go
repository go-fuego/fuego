package cache

import (
	"bytes"
	"log/slog"
	"net/http"
	"time"
)

type Storage interface {
	Get(key string) (string, bool)
	Set(key string, value string)
	Delete(key string)
}

type Config struct {
	Storage Storage
}

// Cache middleware
func New(config Config) func(http.Handler) http.Handler {
	storage := config.Storage

	if storage == nil {
		storage = NewInMemoryCache(3*time.Second, 1000)
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

			key := "httpcache_" + r.URL.Path + "_" + r.Header.Get("Content-Type")

			if val, ok := storage.Get(key); ok {
				respContentType, ok := storage.Get(key + "_response-content-type")
				if ok {
					w.Header().Set("Content-Type", respContentType)
				}

				w.WriteHeader(http.StatusOK)
				slog.Debug("Cache hit", "key", key, "value", string(val), "content-type", respContentType)
				w.Header().Set("Cache", "hit")
				_, _ = w.Write([]byte(val))
				return
			}

			multiWriter := &MultiHTTPWriter{
				ResponseWriter: w,
				cacheWriter:    &bytes.Buffer{},
			}

			h.ServeHTTP(multiWriter, r)

			storage.Set(key, multiWriter.cacheWriter.(*bytes.Buffer).String())
			storage.Set(key+"_response-content-type", multiWriter.Header().Get("Content-Type"))

			slog.Debug("Cached", "key", key)
		})
	}
}
