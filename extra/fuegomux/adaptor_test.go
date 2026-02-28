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

func TestExtractPathParamPatterns(t *testing.T) {
	tests := []struct {
		input    string
		expected map[string]string
	}{
		{"/users/{id:[0-9]+}", map[string]string{"id": "[0-9]+"}},
		{"/users/{id}", map[string]string{}},
		{"/articles/{cat}/{id:[0-9]+}", map[string]string{"id": "[0-9]+"}},
		{"/no-params", map[string]string{}},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			assert.Equal(t, tt.expected, extractPathParamPatterns(tt.input))
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
