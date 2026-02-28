package fuegomux

import (
	"net/http"
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
