package option

import (
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-fuego/fuego"
	"github.com/stretchr/testify/require"
	"github.com/thejerf/slogassert"
)

// dummyMiddleware sets the X-Test header on the request and the X-Test-Response header on the response.
func dummyMiddleware(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.Header.Set("X-Test", "test")
		w.Header().Set("X-Test-Response", "response")
		handler.ServeHTTP(w, r)
	})
}

// orderMiddleware sets the X-Test-Order Header on the request and
// X-Test-Response header on the response. It is
// used to test the order execution of our middleware
func orderMiddleware(s string) func(http.Handler) http.Handler {
	return func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			r.Header.Add("X-Test-Order", s)
			w.Header().Set("X-Test-Response", "response")
			handler.ServeHTTP(w, r)
		})
	}
}

func TestPerRouteMiddleware(t *testing.T) {
	s := fuego.NewServer()

	fuego.Get(s, "/withMiddleware", func(ctx *fuego.ContextNoBody) (string, error) {
		return "withmiddleware", nil
	}, Middleware(dummyMiddleware))

	fuego.Get(s, "/withoutMiddleware", func(ctx *fuego.ContextNoBody) (string, error) {
		return "withoutmiddleware", nil
	})

	t.Run("withMiddleware", func(t *testing.T) {
		r := httptest.NewRequest(http.MethodGet, "/withMiddleware", nil)

		w := httptest.NewRecorder()

		s.Mux.ServeHTTP(w, r)

		require.Equal(t, "withmiddleware", w.Body.String())
		require.Equal(t, "response", w.Header().Get("X-Test-Response"))
	})

	t.Run("withoutMiddleware", func(t *testing.T) {
		r := httptest.NewRequest(http.MethodGet, "/withoutMiddleware", nil)

		w := httptest.NewRecorder()

		s.Mux.ServeHTTP(w, r)

		require.Equal(t, "withoutmiddleware", w.Body.String())
		require.Equal(t, "", w.Header().Get("X-Test-Response"))
	})
}

func TestUse(t *testing.T) {
	t.Run("base", func(t *testing.T) {
		s := fuego.NewServer()
		fuego.Use(s, orderMiddleware("First!"))
		fuego.Get(s, "/test", func(ctx *fuego.ContextNoBody) (string, error) {
			return "test", nil
		})

		r := httptest.NewRequest(http.MethodGet, "/test", nil)
		r.Header.Set("X-Test-Order", "Start!")
		w := httptest.NewRecorder()

		s.Mux.ServeHTTP(w, r)

		require.Equal(t, []string{"Start!", "First!"}, r.Header["X-Test-Order"])
	})

	t.Run("multiple uses of Use", func(t *testing.T) {
		s := fuego.NewServer()
		fuego.Use(s, orderMiddleware("First!"))
		fuego.Use(s, orderMiddleware("Second!"))
		fuego.Get(s, "/test", func(ctx *fuego.ContextNoBody) (string, error) {
			return "test", nil
		})

		r := httptest.NewRequest(http.MethodGet, "/test", nil)
		r.Header.Set("X-Test-Order", "Start!")
		w := httptest.NewRecorder()

		s.Mux.ServeHTTP(w, r)

		require.Equal(t, []string{"Start!", "First!", "Second!"}, r.Header["X-Test-Order"])
	})

	t.Run("variadic use of Use", func(t *testing.T) {
		s := fuego.NewServer()
		fuego.Use(s, orderMiddleware("First!"))
		fuego.Use(s, orderMiddleware("Second!"), orderMiddleware("Third!"))
		fuego.Get(s, "/test", func(ctx *fuego.ContextNoBody) (string, error) {
			return "test", nil
		})

		r := httptest.NewRequest(http.MethodGet, "/test", nil)
		r.Header.Set("X-Test-Order", "Start!")
		w := httptest.NewRecorder()

		s.Mux.ServeHTTP(w, r)

		require.Equal(t, []string{"Start!", "First!", "Second!", "Third!"}, r.Header["X-Test-Order"])
	})

	t.Run("variadic use of Route Get", func(t *testing.T) {
		s := fuego.NewServer()
		fuego.Use(s, orderMiddleware("First!"))
		fuego.Use(s, orderMiddleware("Second!"), orderMiddleware("Third!"))
		fuego.Get(s, "/test", func(ctx *fuego.ContextNoBody) (string, error) {
			return "test", nil
		},
			Middleware(orderMiddleware("Fourth!")),
			Middleware(orderMiddleware("Fifth!")),
		)

		r := httptest.NewRequest(http.MethodGet, "/test", nil)
		r.Header.Set("X-Test-Order", "Start!")
		w := httptest.NewRecorder()

		s.Mux.ServeHTTP(w, r)

		require.Equal(t, []string{"Start!", "First!", "Second!", "Third!", "Fourth!", "Fifth!"}, r.Header["X-Test-Order"])
	})
}

type ans struct{}

func TestParam(t *testing.T) {
	t.Run("warn if param is not found in openAPI config but called in controller (possibly typo)", func(t *testing.T) {
		handler := slogassert.New(t, slog.LevelWarn, nil)

		s := fuego.NewServer(
			fuego.WithLogHandler(handler),
		)

		fuego.Get(s, "/correct", func(c fuego.ContextNoBody) (ans, error) {
			c.QueryParam("quantity")
			return ans{}, nil
		}, Query("quantity", "some description"))

		fuego.Get(s, "/typo", func(c fuego.ContextNoBody) (ans, error) {
			c.QueryParam("quantityy-with-a-typo")
			return ans{}, nil
		}).QueryParam("quantity", "some description")

		t.Run("correct param", func(t *testing.T) {
			r := httptest.NewRequest("GET", "/correct", nil)
			w := httptest.NewRecorder()

			s.Mux.ServeHTTP(w, r)

			// all log messages have been accounted for
			handler.AssertEmpty()
		})

		t.Run("typo param", func(t *testing.T) {
			r := httptest.NewRequest("GET", "/typo", nil)
			w := httptest.NewRecorder()

			s.Mux.ServeHTTP(w, r)

			handler.AssertMessage("query parameter not expected in OpenAPI spec")

			// all log messages have been accounted for
			handler.AssertEmpty()
		})
	})
}
