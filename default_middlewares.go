package fuego

import (
	"log/slog"
	"net/http"
)

type loggingConfig struct {
	Enabled         bool
	DisableRequest  bool
	DisableResponse bool

	RequestLogger  func(w http.ResponseWriter, r *http.Request)
	ResponseLogger func(w http.ResponseWriter, r *http.Request)
}

func defaultLoggingMiddleware(s *Server) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !s.loggingConfig.DisableRequest {
				s.loggingConfig.RequestLogger(w, r)
			}

			next.ServeHTTP(w, r)

			if !s.loggingConfig.DisableResponse {
				s.loggingConfig.ResponseLogger(w, r)
			}
		})
	}
}

func defaultRequestLog(w http.ResponseWriter, r *http.Request) {
	slog.Info("<- request")
}
func defaultResponseLog(w http.ResponseWriter, r *http.Request) {
	slog.Info("response ->")
}
