package fuego

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type LoggingConfig struct {
	DisableRequest  bool // If true, request logging is disabled
	DisableResponse bool // If true, response logging is disabled
}

// responseWriter wraps [http.ResponseWriter] to capture response metadata.
// Implements [http.ResponseWriter.Write] to ensure proper status code capture for implicit 200 responses
type responseWriter struct {
	http.ResponseWriter
	status      int
	wroteHeader bool
}

func wrapResponseWriter(w http.ResponseWriter) *responseWriter {
	return &responseWriter{ResponseWriter: w}
}

func (rw *responseWriter) Status() int {
	return rw.status
}

func (rw *responseWriter) WriteHeader(code int) {
	if rw.wroteHeader {
		return
	}

	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
	rw.wroteHeader = true
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	if !rw.wroteHeader {
		rw.WriteHeader(http.StatusOK)
	}
	return rw.ResponseWriter.Write(b)
}

func (l *LoggingConfig) Disabled() bool {
	return l.DisableRequest && l.DisableResponse
}

// By default, all logging is enabled
var defaultLoggingConfig = LoggingConfig{}

// defaultLoggingMiddleware is this default middleware that logs incoming requests and outgoing responses.
//
// By default, request logging will be logged at the debug level, and response
// logging will be logged at the info level
//
// Log levels managed by [WithLogHandler]
func defaultLoggingMiddleware(s *Server) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			requestID := r.Header.Get("X-Request-ID")
			if requestID == "" {
				requestID = uuid.New().String()
			}
			w.Header().Set("X-Request-ID", requestID)

			wrapped := wrapResponseWriter(w)

			if !s.loggingConfig.DisableRequest {
				logRequest(requestID, r)
			}

			next.ServeHTTP(wrapped, r)

			if !s.loggingConfig.DisableResponse {
				duration := time.Since(start)
				logResponse(r, wrapped, requestID, duration)
			}
		})
	}
}

func logRequest(requestID string, r *http.Request) {
	slog.Debug("incoming request",
		"request_id", requestID,
		"method", r.Method,
		"path", r.URL.Path,
		"timestamp", time.Now().Format(time.RFC3339),
		"remote_addr", r.RemoteAddr,
		"user_agent", r.UserAgent(),
	)
}

func logResponse(r *http.Request, rw *responseWriter, requestID string, duration time.Duration) {
	slog.Info("outgoing response",
		"request_id", requestID,
		"method", r.Method,
		"path", r.URL.Path,
		"timestamp", time.Now().Format(time.RFC3339),
		"duration_ms", duration.Milliseconds(),
		"status_code", rw.status,
	)
}
