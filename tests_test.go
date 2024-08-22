package fuego

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

// Contains random tests reported on the issues.

func TestContentType(t *testing.T) {
	server := NewServer()

	t.Run("Sends application/problem+json when return type is HTTPError", func(t *testing.T) {
		GetStd(server, "/json-problems", func(w http.ResponseWriter, r *http.Request) {
			SendJSONError(w, nil, UnauthorizedError{
				Title: "Unauthorized",
			})
		})

		req := httptest.NewRequest("GET", "/json-problems", nil)
		w := httptest.NewRecorder()
		server.Mux.ServeHTTP(w, req)

		require.Equal(t, "application/problem+json", w.Result().Header.Get("Content-Type"))
		require.Equal(t, 401, w.Code)
		require.Contains(t, w.Body.String(), "Unauthorized")
	})

	t.Run("Sends application/json when return type is not HTTPError", func(t *testing.T) {
		GetStd(server, "/json", func(w http.ResponseWriter, r *http.Request) {
			SendJSONError(w, nil, errors.New("error"))
		})

		req := httptest.NewRequest("GET", "/json", nil)
		w := httptest.NewRecorder()
		server.Mux.ServeHTTP(w, req)

		require.Equal(t, "application/json", w.Header().Get("Content-Type"))
		require.Equal(t, 500, w.Code)
		require.Equal(t, "{}\n", w.Body.String())
	})
}
