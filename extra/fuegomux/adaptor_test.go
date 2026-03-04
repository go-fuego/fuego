package fuegomux

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"

	"github.com/go-fuego/fuego"
)

func TestMuxToFuegoRoute(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"/users/{id}", "/users/{id}"},
		{"/users/{id:[0-9]+}", "/users/{id}"},
		{"/articles/{category}/{id:[0-9]+}", "/articles/{category}/{id}"},
		{"/files/{path:.*}", "/files/{path}"},
		{"/sort/{order:(?:asc|desc)}", "/sort/{order}"},
		{"/no-params", "/no-params"},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			assert.Equal(t, tt.expected, muxToFuegoRoute(tt.input))
		})
	}
}

func TestFuegoRouteRegistration(t *testing.T) {
	e := fuego.NewEngine()
	r := mux.NewRouter()

	Get(e, r, "/users", func(c fuego.ContextNoBody) (string, error) { return "ok", nil })

	spec := e.OutputOpenAPISpec()
	assert.NotNil(t, spec.Paths.Find("/users"))
}

func TestFuegoRouteWithPathParam(t *testing.T) {
	e := fuego.NewEngine()
	r := mux.NewRouter()

	Get(e, r, "/users/{id}", func(c fuego.ContextNoBody) (string, error) { return "ok", nil })

	spec := e.OutputOpenAPISpec()
	assert.NotNil(t, spec.Paths.Find("/users/{id}"))
}

func TestFuegoRouteWithRegexParam(t *testing.T) {
	e := fuego.NewEngine()
	r := mux.NewRouter()

	Get(e, r, "/users/{id:[0-9]+}", func(c fuego.ContextNoBody) (string, error) { return "ok", nil })

	spec := e.OutputOpenAPISpec()
	// Regex should be stripped from the OpenAPI path
	assert.NotNil(t, spec.Paths.Find("/users/{id}"))
}

func TestFuegoRouteWithSubrouter(t *testing.T) {
	e := fuego.NewEngine()
	r := mux.NewRouter()
	sub := r.PathPrefix("/api").Subrouter()

	Get(e, sub, "/users", func(c fuego.ContextNoBody) (string, error) { return "ok", nil })

	spec := e.OutputOpenAPISpec()
	assert.NotNil(t, spec.Paths.Find("/api/users"))
}

func TestMuxHandlerRegistration(t *testing.T) {
	e := fuego.NewEngine()
	r := mux.NewRouter()

	GetMux(e, r, "/native", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	spec := e.OutputOpenAPISpec()
	assert.NotNil(t, spec.Paths.Find("/native"))
}

func TestFuegoHandler_Integration(t *testing.T) {
	e := fuego.NewEngine()
	r := mux.NewRouter()

	type Response struct {
		ID      string `json:"id"`
		Message string `json:"message"`
	}

	Get(e, r, "/users/{id:[0-9]+}", func(c fuego.ContextNoBody) (Response, error) {
		return Response{
			ID:      c.PathParam("id"),
			Message: "hello",
		}, nil
	})

	req := httptest.NewRequest(http.MethodGet, "/users/42", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), `"id":"42"`)
	assert.Contains(t, w.Body.String(), `"message":"hello"`)
}

func TestFuegoHandler_PostBody(t *testing.T) {
	e := fuego.NewEngine()
	r := mux.NewRouter()

	type Request struct {
		Name string `json:"name" validate:"required"`
	}
	type Response struct {
		Greeting string `json:"greeting"`
	}

	Post(e, r, "/greet", func(c fuego.ContextWithBody[Request]) (Response, error) {
		body, err := c.Body()
		if err != nil {
			return Response{}, err
		}
		return Response{Greeting: "Hello " + body.Name}, nil
	})

	req := httptest.NewRequest(http.MethodPost, "/greet", strings.NewReader(`{"name":"World"}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), `"greeting":"Hello World"`)
}

func TestOptionMiddleware_Applied(t *testing.T) {
	e := fuego.NewEngine()
	r := mux.NewRouter()

	middlewareCalled := false
	testMiddleware := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			middlewareCalled = true
			next.ServeHTTP(w, r)
		})
	}

	Get(e, r, "/protected", func(c fuego.ContextNoBody) (string, error) {
		return "ok", nil
	}, fuego.OptionMiddleware(testMiddleware))

	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.True(t, middlewareCalled, "OptionMiddleware should be applied to the handler")
}

func TestOptionMiddleware_CanBlockRequest(t *testing.T) {
	e := fuego.NewEngine()
	r := mux.NewRouter()

	authMiddleware := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Header.Get("Authorization") == "" {
				http.Error(w, "unauthorized", http.StatusForbidden)
				return
			}
			next.ServeHTTP(w, r)
		})
	}

	handlerCalled := false
	Get(e, r, "/secret", func(c fuego.ContextNoBody) (string, error) {
		handlerCalled = true
		return "secret data", nil
	}, fuego.OptionMiddleware(authMiddleware))

	// Request without auth header — should be blocked
	req := httptest.NewRequest(http.MethodGet, "/secret", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusForbidden, w.Code)
	assert.False(t, handlerCalled, "handler should not be called when middleware blocks")

	// Request with auth header — should pass
	handlerCalled = false
	req = httptest.NewRequest(http.MethodGet, "/secret", nil)
	req.Header.Set("Authorization", "Bearer token")
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.True(t, handlerCalled, "handler should be called when middleware passes")
}

func TestOptionMiddleware_AppliedToMuxHandler(t *testing.T) {
	e := fuego.NewEngine()
	r := mux.NewRouter()

	middlewareCalled := false
	testMiddleware := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			middlewareCalled = true
			next.ServeHTTP(w, r)
		})
	}

	GetMux(e, r, "/native-protected", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}, fuego.OptionMiddleware(testMiddleware))

	req := httptest.NewRequest(http.MethodGet, "/native-protected", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.True(t, middlewareCalled, "OptionMiddleware should be applied to native mux handlers too")
}

func TestOptionMiddleware_MultipleMiddlewares_Order(t *testing.T) {
	e := fuego.NewEngine()
	r := mux.NewRouter()

	var order []string
	mw1 := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			order = append(order, "mw1-before")
			next.ServeHTTP(w, r)
			order = append(order, "mw1-after")
		})
	}
	mw2 := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			order = append(order, "mw2-before")
			next.ServeHTTP(w, r)
			order = append(order, "mw2-after")
		})
	}

	Get(e, r, "/ordered", func(c fuego.ContextNoBody) (string, error) {
		order = append(order, "handler")
		return "ok", nil
	}, fuego.OptionMiddleware(mw1, mw2))

	req := httptest.NewRequest(http.MethodGet, "/ordered", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, []string{"mw1-before", "mw2-before", "handler", "mw2-after", "mw1-after"}, order)
}

func TestBody_DisallowUnknownFields(t *testing.T) {
	e := fuego.NewEngine()
	r := mux.NewRouter()

	type Request struct {
		Name string `json:"name"`
	}

	Post(e, r, "/strict", func(c fuego.ContextWithBody[Request]) (string, error) {
		_, err := c.Body()
		if err != nil {
			return "", err
		}
		return "ok", nil
	})

	// Request with unknown field — should be rejected
	req := httptest.NewRequest(http.MethodPost, "/strict", strings.NewReader(`{"name":"test","unknown":"field"}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.NotEqual(t, http.StatusOK, w.Code, "unknown fields should be rejected")
}

func TestBody_XML(t *testing.T) {
	e := fuego.NewEngine()
	r := mux.NewRouter()

	type Request struct {
		Name string `xml:"name"`
	}
	type Response struct {
		Greeting string `json:"greeting"`
	}

	Post(e, r, "/greet", func(c fuego.ContextWithBody[Request]) (Response, error) {
		body, err := c.Body()
		if err != nil {
			return Response{}, err
		}
		return Response{Greeting: "Hello " + body.Name}, nil
	})

	req := httptest.NewRequest(http.MethodPost, "/greet", strings.NewReader(`<Request><name>World</name></Request>`))
	req.Header.Set("Content-Type", "application/xml")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), `"greeting":"Hello World"`)
}

func TestBody_URLEncoded(t *testing.T) {
	e := fuego.NewEngine()
	r := mux.NewRouter()

	type Request struct {
		Name string `schema:"name"`
	}
	type Response struct {
		Greeting string `json:"greeting"`
	}

	Post(e, r, "/greet", func(c fuego.ContextWithBody[Request]) (Response, error) {
		body, err := c.Body()
		if err != nil {
			return Response{}, err
		}
		return Response{Greeting: "Hello " + body.Name}, nil
	})

	req := httptest.NewRequest(http.MethodPost, "/greet", strings.NewReader("name=World"))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), `"greeting":"Hello World"`)
}

func TestBody_DefaultsToJSON(t *testing.T) {
	e := fuego.NewEngine()
	r := mux.NewRouter()

	type Request struct {
		Name string `json:"name"`
	}

	Post(e, r, "/test", func(c fuego.ContextWithBody[Request]) (Request, error) {
		return c.Body()
	})

	// No Content-Type header — should default to JSON
	req := httptest.NewRequest(http.MethodPost, "/test", strings.NewReader(`{"name":"test"}`))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), `"name":"test"`)
}

func TestSubrouterMiddleware_WithOptionMiddleware(t *testing.T) {
	e := fuego.NewEngine()
	r := mux.NewRouter()

	var order []string

	// Subrouter with group-level middleware (like PortalSessionRequiredHandler)
	protected := r.PathPrefix("").Subrouter()
	protected.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			order = append(order, "group-middleware")
			next.ServeHTTP(w, r)
		})
	})

	// Route with per-route middleware (like PortalAuthHandler)
	routeMiddleware := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			order = append(order, "route-middleware")
			next.ServeHTTP(w, r)
		})
	}

	Get(e, protected, "/resource", func(c fuego.ContextNoBody) (string, error) {
		order = append(order, "handler")
		return "ok", nil
	}, fuego.OptionMiddleware(routeMiddleware))

	req := httptest.NewRequest(http.MethodGet, "/resource", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	// Group middleware runs first (applied by gorilla/mux subrouter),
	// then route-level middleware (applied by fuegomux), then handler
	assert.Equal(t, []string{"group-middleware", "route-middleware", "handler"}, order)
}
