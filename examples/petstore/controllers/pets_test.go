package controller_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/go-fuego/fuego/examples/petstore/lib"
)

func TestGetAllPets(t *testing.T) {
	t.Run("can get all pets", func(t *testing.T) {
		s := lib.NewPetStoreServer()

		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/pets/all?per_page=5", nil)

		s.Mux.ServeHTTP(w, r)

		require.Equal(t, http.StatusOK, w.Code)
	})
}

func TestFilterPets(t *testing.T) {
	t.Run("can filter pets", func(t *testing.T) {
		s := lib.NewPetStoreServer()

		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/pets/?name=kit&per_page=5", nil)
		s.Mux.ServeHTTP(w, r)
		t.Log(w.Body.String())
		require.Equal(t, http.StatusOK, w.Code)

		w = httptest.NewRecorder()
		r = httptest.NewRequest("GET", "/pets/?name=kit&younger_than=1&per_page=5", nil)
		s.Mux.ServeHTTP(w, r)
		t.Log(w.Body.String())
		require.Equal(t, http.StatusOK, w.Code)
	})
}

func TestPostPets(t *testing.T) {
	t.Run("can create a pet", func(t *testing.T) {
		s := lib.NewPetStoreServer()

		w := httptest.NewRecorder()
		r := httptest.NewRequest("POST", "/pets/", strings.NewReader(`{"name": "kitkat"}`))

		s.Mux.ServeHTTP(w, r)

		require.Equal(t, http.StatusOK, w.Code)
	})
}

func TestGetPets(t *testing.T) {
	t.Run("can get a pet by id", func(t *testing.T) {
		s := lib.NewPetStoreServer()

		w := httptest.NewRecorder()
		r := httptest.NewRequest("POST", "/pets/", strings.NewReader(`{"name": "kitkat"}`))
		s.Mux.ServeHTTP(w, r)
		t.Log(w.Body.String())
		require.Equal(t, http.StatusOK, w.Code)

		w = httptest.NewRecorder()
		r = httptest.NewRequest("GET", "/pets/pet-1", nil)

		s.Mux.ServeHTTP(w, r)

		t.Log(w.Body.String())
		require.Equal(t, http.StatusOK, w.Code)
	})
}

func TestGetAllPestByAge(t *testing.T) {
	t.Run("can get a pet by id", func(t *testing.T) {
		s := lib.NewPetStoreServer()

		w := httptest.NewRecorder()
		r := httptest.NewRequest("POST", "/pets/", strings.NewReader(`{"name": "kitkat"}`))
		s.Mux.ServeHTTP(w, r)
		t.Log(w.Body.String())
		require.Equal(t, http.StatusOK, w.Code)

		w = httptest.NewRecorder()
		r = httptest.NewRequest("GET", "/pets/by-age", nil)

		s.Mux.ServeHTTP(w, r)

		t.Log(w.Body.String())
		require.Equal(t, http.StatusOK, w.Code)
	})
}

func TestGetPetsByName(t *testing.T) {
	t.Run("can get a pet by name", func(t *testing.T) {
		s := lib.NewPetStoreServer()

		w := httptest.NewRecorder()
		r := httptest.NewRequest("POST", "/pets/", strings.NewReader(`{"name": "kitkat"}`))
		s.Mux.ServeHTTP(w, r)
		t.Log(w.Body.String())
		require.Equal(t, http.StatusOK, w.Code)

		w = httptest.NewRecorder()
		r = httptest.NewRequest("GET", "/pets/by-name/kitkat", nil)

		s.Mux.ServeHTTP(w, r)

		t.Log(w.Body.String())
		require.Equal(t, http.StatusOK, w.Code)
	})
}

func TestPutPets(t *testing.T) {
	t.Run("can update a pet", func(t *testing.T) {
		s := lib.NewPetStoreServer()

		w := httptest.NewRecorder()
		r := httptest.NewRequest("POST", "/pets/", strings.NewReader(`{"name": "kitkat"}`))
		s.Mux.ServeHTTP(w, r)
		t.Log(w.Body.String())
		require.Equal(t, http.StatusOK, w.Code)

		w = httptest.NewRecorder()
		r = httptest.NewRequest("PUT", "/pets/pet-1", strings.NewReader(`{"name": "snickers"}`))

		s.Mux.ServeHTTP(w, r)

		t.Log(w.Body.String())
		require.Equal(t, http.StatusOK, w.Code)
	})
}

func TestDeletePets(t *testing.T) {
	t.Run("can delete a pet", func(t *testing.T) {
		s := lib.NewPetStoreServer()

		w := httptest.NewRecorder()
		r := httptest.NewRequest("POST", "/pets/", strings.NewReader(`{"name": "kitkat"}`))
		s.Mux.ServeHTTP(w, r)
		t.Log(w.Body.String())
		require.Equal(t, http.StatusOK, w.Code)

		w = httptest.NewRecorder()
		r = httptest.NewRequest("DELETE", "/pets/pet-1", nil)

		s.Mux.ServeHTTP(w, r)

		t.Log(w.Body.String())
		require.Equal(t, http.StatusOK, w.Code)
	})
}
