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
			if !s.loggingConfig.Enabled {
				next.ServeHTTP(w, r)
				return
			}

			// Placeholders (will flesh this out later)
			if !s.loggingConfig.DisableRequest {
				if s.loggingConfig.RequestLogger != nil {
					s.loggingConfig.RequestLogger(w, r)
				} else {
					slog.Info("<- request")
				}
			}

			next.ServeHTTP(w, r)

			if !s.loggingConfig.DisableResponse {
				if s.loggingConfig.ResponseLogger != nil {
					s.loggingConfig.ResponseLogger(w, r)
				} else {
					slog.Info("response ->")
				}
			}
		})
	}

}
