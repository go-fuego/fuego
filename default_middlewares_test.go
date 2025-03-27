package fuego

import (
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/thejerf/slogassert"
)

type TestResponseWriter struct{}

func (w *TestResponseWriter) Header() http.Header {
	panic("not implemented")
}

func (w *TestResponseWriter) Write(b []byte) (int, error) {
	panic("not implemented")
}

func (w *TestResponseWriter) WriteHeader(statusCode int) {
	panic("not implemented")
}

func TestFlush(t *testing.T) {
	t.Run("is implemented", func(t *testing.T) {
		handler := slogassert.New(t, slog.LevelWarn, nil)
		slog.SetDefault(slog.New(handler))

		rw := newResponseWriter(httptest.NewRecorder())
		rw.Flush()
		handler.AssertEmpty()
	})
	t.Run("is not implemented", func(t *testing.T) {
		handler := slogassert.New(t, slog.LevelWarn, nil)
		slog.SetDefault(slog.New(handler))

		rw := newResponseWriter(&TestResponseWriter{})
		rw.Flush()
		handler.AssertMessage("Flush not implemented, skipping")
		handler.AssertEmpty()
	})
}
