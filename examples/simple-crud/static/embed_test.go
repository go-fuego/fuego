package static

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandler(t *testing.T) {
	h := Handler()

	t.Run("200 with raw handler", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/tailwind.min.css", nil)
		h.ServeHTTP(w, r)

		if w.Code != 200 {
			t.Errorf("Handler() = %v, want %v", w.Code, 200)
		}
	})

	t.Run("404 with raw handler", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/not-existing", nil)
		h.ServeHTTP(w, r)

		if w.Code != 404 {
			t.Errorf("Handler() = %v, want %v", w.Code, 404)
		}
	})

	t.Run("200 with mux handler", func(t *testing.T) {
		mux := http.NewServeMux()
		mux.Handle("/static/", http.StripPrefix("/static", h))

		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/static/tailwind.min.css", nil)
		mux.ServeHTTP(w, r)

		if w.Code != 200 {
			t.Errorf("Handler() = %v, want %v", w.Code, 200)
		}

		t.Log(w.Body.String())
		t.Fail()
	})
}
