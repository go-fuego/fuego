package httpLog

import (
	"log/slog"
	"net/http"
	"time"
)

func New() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			method := r.Method
			path := r.URL.Path

			start := time.Now()
			slog.Info("Received " + method + " " + path)
			next.ServeHTTP(w, r)
			slog.Info("Responded "+method+" "+path, "in", time.Since(start))
		})
	}
}
