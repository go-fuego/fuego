package fuego_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/go-fuego/fuego"
	"github.com/go-fuego/fuego/option"
	"github.com/go-fuego/fuego/param"
)

func TestParamsValidation(t *testing.T) {
	t.Run("Should enforce Required query parameter", func(t *testing.T) {
		s := fuego.NewServer()

		fuego.Get(s, "/test", dummyController,
			option.Query("name", "Name", param.Required(), param.Example("example1", "you")),
		)
		r := httptest.NewRequest("GET", "/test", nil)
		w := httptest.NewRecorder()
		s.Mux.ServeHTTP(w, r)
		require.Equal(t, http.StatusBadRequest, w.Code)
		require.Contains(t, w.Body.String(), "name is a required query param")
	})

	t.Run("Should enforce Required header", func(t *testing.T) {
		s := fuego.NewServer()

		fuego.Get(s, "/test", dummyController,
			option.Header("foo", "header that is foo", param.Required()),
		)
		r := httptest.NewRequest("GET", "/test", nil)
		w := httptest.NewRecorder()
		s.Mux.ServeHTTP(w, r)
		require.Equal(t, http.StatusBadRequest, w.Code)
		require.Contains(t, w.Body.String(), "foo is a required header")
	})

	t.Run("Should enforce Required cookie", func(t *testing.T) {
		s := fuego.NewServer()

		fuego.Get(s, "/test", dummyController,
			option.Cookie("bar", "cookie that is bar", param.Required()),
		)
		r := httptest.NewRequest("GET", "/test", nil)
		w := httptest.NewRecorder()
		s.Mux.ServeHTTP(w, r)
		require.Equal(t, http.StatusBadRequest, w.Code)
		require.Contains(t, w.Body.String(), "bar is a required cookie")
	})
}
