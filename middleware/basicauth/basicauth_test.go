package basicauth_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-fuego/fuego/middleware/basicauth"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	t.Run("cannot create middleware without username", func(t *testing.T) {
		require.Panics(t, func() {
			basicauth.New(basicauth.Config{
				Password: "pass",
			})
		})
	})

	t.Run("cannot create middleware without password", func(t *testing.T) {
		require.Panics(t, func() {
			basicauth.New(basicauth.Config{
				Username: "user",
			})
		})
	})

	t.Run("usual creation", func(t *testing.T) {
		basicAuth := basicauth.New(basicauth.Config{
			Username: "user",
			Password: "pass",
		})
		require.NotNil(t, basicAuth)

		handler := basicAuth(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("OK"))
		}))

		t.Run("without auth", func(t *testing.T) {
			req, err := http.NewRequest(http.MethodGet, "/", nil)
			require.NoError(t, err)

			// Test without auth
			w := httptest.NewRecorder()
			handler.ServeHTTP(w, req)
			require.Equal(t, http.StatusUnauthorized, w.Code)
			require.Equal(t, `Basic realm="Restricted"`, w.Header().Get("WWW-Authenticate"))
		})

		t.Run("with wrong auth", func(t *testing.T) {
			req, err := http.NewRequest(http.MethodGet, "/", nil)
			require.NoError(t, err)
			req.SetBasicAuth("user", "wrongpass")

			// Test with wrong auth
			w := httptest.NewRecorder()
			handler.ServeHTTP(w, req)
			require.Equal(t, http.StatusUnauthorized, w.Code)
			require.Equal(t, `Basic realm="Restricted"`, w.Header().Get("WWW-Authenticate"))
		})

		t.Run("with correct auth", func(t *testing.T) {
			req, err := http.NewRequest(http.MethodGet, "/", nil)
			require.NoError(t, err)
			req.SetBasicAuth("user", "pass")

			// Test with correct auth
			w := httptest.NewRecorder()
			handler.ServeHTTP(w, req)
			require.Equal(t, http.StatusOK, w.Code)
			require.Equal(t, "OK", w.Body.String())
		})
	})

	t.Run("allow get requests without auth", func(t *testing.T) {
		basicAuth := basicauth.New(basicauth.Config{
			Username: "user",
			Password: "pass",
			AllowGet: true,
		})
		require.NotNil(t, basicAuth)

		handler := basicAuth(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("OK"))
		}))

		t.Run("get authorized without auth", func(t *testing.T) {
			req, err := http.NewRequest(http.MethodGet, "/", nil)
			require.NoError(t, err)

			// Test without auth
			w := httptest.NewRecorder()
			handler.ServeHTTP(w, req)
			require.Equal(t, http.StatusOK, w.Code)
			require.Equal(t, "OK", w.Body.String())
		})

		t.Run("get authorized with wrong auth", func(t *testing.T) {
			req, err := http.NewRequest(http.MethodGet, "/", nil)
			require.NoError(t, err)
			req.SetBasicAuth("user", "wrongpass")

			// Test with wrong auth
			w := httptest.NewRecorder()
			handler.ServeHTTP(w, req)
			require.Equal(t, http.StatusOK, w.Code)
			require.Equal(t, "OK", w.Body.String())
		})

		t.Run("post unauthorized without auth", func(t *testing.T) {
			req, err := http.NewRequest(http.MethodPost, "/", nil)
			require.NoError(t, err)

			// Test without auth
			w := httptest.NewRecorder()
			handler.ServeHTTP(w, req)
			require.Equal(t, http.StatusUnauthorized, w.Code)
			require.Equal(t, `Basic realm="Restricted"`, w.Header().Get("WWW-Authenticate"))
		})
	})
}
