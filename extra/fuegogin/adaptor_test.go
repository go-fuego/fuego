package fuegogin_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/go-fuego/fuego/extra/fuegogin/lib"
)

func TestFuegoGin(t *testing.T) {
	e, _ := lib.SetupGin()

	t.Run("simply test gin", func(t *testing.T) {
		r := httptest.NewRequest("GET", "/gin", nil)
		w := httptest.NewRecorder()

		e.ServeHTTP(w, r)

		require.Equal(t, 200, w.Code)
	})

	t.Run("test fuego plugin", func(t *testing.T) {
		r := httptest.NewRequest("GET", "/fuego", nil)
		w := httptest.NewRecorder()

		e.ServeHTTP(w, r)

		require.Equal(t, http.StatusOK, w.Code)
		require.JSONEq(t, `{"message":"Hello "}`, w.Body.String())
	})
}
