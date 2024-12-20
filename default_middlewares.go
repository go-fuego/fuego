package fuego

import (
	"log/slog"
	"net/http"
)

func (l *RequestResponseLogger) Log(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		l.request(w, r)
		next.ServeHTTP(w, r)
		l.response(w, r)
	})
}

func RequestLog(w http.ResponseWriter, r *http.Request) {
	slog.Info("request")
}

func ResponseLog(w http.ResponseWriter, r *http.Request) {
	slog.Info("response")
}
