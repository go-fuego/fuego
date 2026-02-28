package fuegomux

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"

	"github.com/go-fuego/fuego"
	"github.com/go-fuego/fuego/internal"
)

// Compile-time interface checks
var (
	_ fuego.ContextFlowable[any, any] = &muxContext[any, any]{}
	_ fuego.Context[any, any]         = &muxContext[any, any]{}
	_ fuego.ContextWithBody[any]      = &muxContext[any, any]{}
)

func TestMuxContext_PathParam(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := &muxContext[any, any]{
			CommonContext: internal.CommonContext[any]{
				CommonCtx: r.Context(),
				UrlValues: r.URL.Query(),
			},
			req: r,
			res: w,
		}
		assert.Equal(t, "42", ctx.PathParam("id"))
		assert.Equal(t, 42, ctx.PathParamInt("id"))
		w.WriteHeader(http.StatusOK)
	})

	r := mux.NewRouter()
	r.HandleFunc("/users/{id}", handler)

	req := httptest.NewRequest(http.MethodGet, "/users/42", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestMuxContext_RequestResponse(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := &muxContext[any, any]{
			CommonContext: internal.CommonContext[any]{
				CommonCtx: r.Context(),
				UrlValues: r.URL.Query(),
			},
			req: r,
			res: w,
		}
		assert.Equal(t, r, ctx.Request())
		assert.Equal(t, w, ctx.Response())
		assert.Equal(t, "bar", ctx.Header("X-Foo"))
		ctx.SetHeader("X-Out", "baz")
		assert.Equal(t, "baz", w.Header().Get("X-Out"))
		w.WriteHeader(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("X-Foo", "bar")
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)
}
