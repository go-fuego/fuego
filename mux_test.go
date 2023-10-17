package op

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func dummyMiddleware(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.Header.Set("X-Test", "test")
		handler.ServeHTTP(w, r)
	})
}

func TestUseStd(t *testing.T) {
	s := NewServer()
	UseStd(s, dummyMiddleware)
	GetStd(s, "/test", func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-Test") != "test" {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("middleware not registered"))
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("test successfull"))
	})

	r := httptest.NewRequest(http.MethodGet, "/test", nil)

	w := httptest.NewRecorder()

	s.mux.ServeHTTP(w, r)

	require.Equal(t, w.Code, http.StatusOK)
	require.Equal(t, w.Body.String(), "test successfull")
}
