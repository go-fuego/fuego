package fuego

import (
	"log/slog"
	"net/http"
)

type LoggingConfig struct {
	Enabled         bool
	DisableRequest  bool
	DisableResponse bool

	RequestLogger  func(w http.ResponseWriter, r *http.Request)
	ResponseLogger func(w http.ResponseWriter, r *http.Request)
}

func defaultLoggingMiddleware(s *Server) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !s.LoggingConfig.Enabled {
				next.ServeHTTP(w, r)
				return
			}

			// Placeholders (will flesh this out later)
			if !s.LoggingConfig.DisableRequest {
				if s.LoggingConfig.RequestLogger != nil {
					s.LoggingConfig.RequestLogger(w, r)
				} else {
					slog.Info("<- request")
				}
			}

			next.ServeHTTP(w, r)

			if !s.LoggingConfig.DisableResponse {
				if s.LoggingConfig.ResponseLogger != nil {
					s.LoggingConfig.ResponseLogger(w, r)
				} else {
					slog.Info("response ->")
				}
			}
		})
	}

}
