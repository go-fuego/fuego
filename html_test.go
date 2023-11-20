package op

import (
	"embed"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

//go:embed testdata/*.html
var testdata embed.FS

func TestRender(t *testing.T) {
	s := NewServer(
		WithTemplateFS(testdata),
		WithTemplateGlobs("testdata/*.html"),
	)

	Get(s, "/test", func(ctx Ctx[any]) (HTML, error) {
		return ctx.Render("testdata/test.html", H{"Name": "test"})
	})

	t.Run("Render once", func(t *testing.T) {
		r := httptest.NewRequest(http.MethodGet, "/test", nil)
		w := httptest.NewRecorder()

		s.mux.ServeHTTP(w, r)

		require.Equal(t, w.Code, http.StatusOK)
		require.Equal(t, w.Body.String(), "<main>\n  <h1>Test</h1>\n  <p>Your name is: test</p>\n</main>\n")
	})

	t.Run("Render twice", func(t *testing.T) {
		r := httptest.NewRequest(http.MethodGet, "/test", nil)
		w := httptest.NewRecorder()

		s.mux.ServeHTTP(w, r)

		require.Equal(t, w.Code, http.StatusOK)
		require.Equal(t, w.Body.String(), "<main>\n  <h1>Test</h1>\n  <p>Your name is: test</p>\n</main>\n")
	})
}

func BenchmarkRender(b *testing.B) {
	s := NewServer(
		WithTemplateFS(testdata),
		WithTemplateGlobs("testdata/*.html"),
	)

	Get(s, "/test", func(ctx Ctx[any]) (HTML, error) {
		return ctx.Render("testdata/test.html", H{"Name": "test"})
	})

	expected := "<main>\n  <h1>Test</h1>\n  <p>Your name is: test</p>\n</main>\n"

	for i := 0; i < b.N; i++ {
		r := httptest.NewRequest(http.MethodGet, "/test", nil)
		w := httptest.NewRecorder()

		s.mux.ServeHTTP(w, r)
		if w.Code != http.StatusOK {
			b.Fail()
		}
		if w.Body.String() != expected {
			b.Fail()
		}
	}
}
