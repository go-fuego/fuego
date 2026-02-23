package testingfromoutside_test

import (
	"log"
	"net/http/httptest"
	"testing"

	"github.com/jub0bs/cors"
	"github.com/stretchr/testify/require"

	"github.com/go-fuego/fuego"
)

func TestCors(t *testing.T) {
	cors, err := cors.NewMiddleware(cors.Config{
		Origins: []string{"*"},
		Methods: []string{"GET"},
	})
	if err != nil {
		log.Fatal(err)
	}

	s := fuego.NewServer(
		fuego.WithoutLogger(),
		fuego.WithGlobalMiddlewares(cors.Wrap),
	)

	fuego.Get(s, "/", func(c fuego.ContextNoBody) (string, error) {
		return "Hello, World!", nil
	})

	t.Run("CORS request INCOMPLETE TEST", func(t *testing.T) {
		r := httptest.NewRequest("GET", "http://example.com/", nil)
		w := httptest.NewRecorder()

		r.Header.Set("Origin", "http://example.com/")
		r.Header.Set("Access-Control-Request-Method", "GET")

		s.Mux.ServeHTTP(w, r)

		t.Log(w.Header())
		body := w.Body.String()
		require.Equal(t, "Hello, World!", body)
		require.Equal(t, 200, w.Code)
	})
}
