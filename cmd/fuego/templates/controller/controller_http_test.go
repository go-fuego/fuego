package controller

import (
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-fuego/fuego"
)

// HTTP tests for /NewController
func TestNewControllerRessources_HttpRoutes(t *testing.T) {
	s := fuego.NewServer()

	rs := NewControllerRessources{
		// Dependency Injection
		NewControllerService: NewControllerServiceMock{
			getAllNewControllerLength: 2,
		},
	}

	rs.Routes(s)

	t.Run("404", func(t *testing.T) {
		r := httptest.NewRequest("GET", "/non-existing", nil)
		w := httptest.NewRecorder()

		s.Mux.ServeHTTP(w, r)

		if w.Code != 404 {
			t.Errorf("Expected status code 404, got %d", w.Code)
		}
	})

	t.Run("GET /newController", func(t *testing.T) {
		r := httptest.NewRequest("GET", "/newController/", nil)
		w := httptest.NewRecorder()

		s.Mux.ServeHTTP(w, r)

		if w.Code != 200 {
			t.Errorf("Expected status code 200, got %d", w.Code)
		}
		expectedBody := `[{"id":"","name":""},{"id":"","name":""}]`
		actualBody := strings.TrimSpace(w.Body.String())
		if expectedBody != actualBody {
			t.Errorf("Expected body %s, got %s", expectedBody, actualBody)
		}
	})

	t.Run("GET /newController/{id}", func(t *testing.T) {
		r := httptest.NewRequest("GET", "/newController/123", nil)
		w := httptest.NewRecorder()

		s.Mux.ServeHTTP(w, r)

		if w.Code != 200 {
			t.Errorf("Expected status code 200, got %d", w.Code)
		}
		expectedBody := `{"id":"","name":""}`
		actualBody := strings.TrimSpace(w.Body.String())
		if expectedBody != actualBody {
			t.Errorf("Expected body %s, got %s", expectedBody, actualBody)
		}
	})
}
