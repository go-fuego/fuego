package cache

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-fuego/fuego"
	"github.com/stretchr/testify/require"
)

type testStruct struct {
	Name string
	Age  int
}

const waitTime = 10 * time.Millisecond

func baseController(ctx fuego.Ctx[any]) (testStruct, error) {
	time.Sleep(waitTime)
	return testStruct{Name: "test", Age: 10}, nil
}

func TestCache(t *testing.T) {
	s := fuego.NewServer()

	fuego.Get(s, "/without-cache", baseController)

	fuego.Use(s, New(Config{}))

	fuego.Get(s, "/with-cache", baseController)

	t.Run("Answer once", func(t *testing.T) {
		r := httptest.NewRequest("GET", "/without-cache", nil)
		w := httptest.NewRecorder()

		s.Mux.ServeHTTP(w, r)

		require.Equal(t, http.StatusOK, w.Code)
		require.Equal(t, w.Body.String(), `{"Name":"test","Age":10}`+"\n")
	})

	t.Run("Answer twice without cache", func(t *testing.T) {
		r := httptest.NewRequest("GET", "/without-cache", nil)
		w := httptest.NewRecorder()

		s.Mux.ServeHTTP(w, r)

		require.Equal(t, http.StatusOK, w.Code)
		require.Equal(t, w.Body.String(), `{"Name":"test","Age":10}`+"\n")

		r = httptest.NewRequest("GET", "/without-cache", nil)
		w = httptest.NewRecorder()

		start := time.Now()
		s.Mux.ServeHTTP(w, r)
		elapsed := time.Since(start)
		require.True(t, elapsed >= waitTime)

		require.Equal(t, http.StatusOK, w.Code)
		require.Equal(t, w.Body.String(), `{"Name":"test","Age":10}`+"\n")
	})

	t.Run("Answer twice with cache", func(t *testing.T) {
		r := httptest.NewRequest("GET", "/with-cache", nil)
		w := httptest.NewRecorder()

		s.Mux.ServeHTTP(w, r)

		require.Equal(t, http.StatusOK, w.Code)
		require.Equal(t, w.Body.String(), `{"Name":"test","Age":10}`+"\n")

		r = httptest.NewRequest("GET", "/with-cache", nil)
		w = httptest.NewRecorder()

		start := time.Now()
		s.Mux.ServeHTTP(w, r)
		elapsed := time.Since(start)
		require.True(t, elapsed < waitTime)

		require.Equal(t, http.StatusOK, w.Code)
		require.Equal(t, w.Body.String(), `{"Name":"test","Age":10}`+"\n")
	})
}

func BenchmarkCache(b *testing.B) {
	s := fuego.NewServer()

	fuego.Get(s, "/without-cache", baseController)

	fuego.Use(s, New(Config{}))

	fuego.Get(s, "/with-cache", baseController)

	b.Run("without cache", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			r := httptest.NewRequest("GET", "/without-cache", nil)
			w := httptest.NewRecorder()

			s.Mux.ServeHTTP(w, r)

			if w.Code != http.StatusOK {
				b.Fail()
			}
			require.Equal(b, w.Body.String(), `{"Name":"test","Age":10}`+"\n")
		}
	})

	b.Run("with cache", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			r := httptest.NewRequest("GET", "/with-cache", nil)
			w := httptest.NewRecorder()

			s.Mux.ServeHTTP(w, r)

			if w.Code != http.StatusOK {
				b.Fail()
			}
			require.Equal(b, w.Body.String(), `{"Name":"test","Age":10}`+"\n")
		}
	})
}
