package fuegomux

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/go-fuego/fuego"
	"github.com/go-fuego/fuego/internal"
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

func TestMuxContext_BodyUnsupportedMediaType(t *testing.T) {
	tests := []struct {
		name        string
		contentType string
		body        string
	}{
		{name: "text/plain", contentType: "text/plain", body: "Hello World"},
		{name: "application/octet-stream", contentType: "application/octet-stream", body: "image"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(tt.body))
			req.Header.Set("Content-Type", tt.contentType)

			ctx := &muxContext[struct{ Name string }, any]{
				CommonContext: internal.CommonContext[struct{ Name string }]{
					CommonCtx: req.Context(),
					UrlValues: req.URL.Query(),
				},
				req: req,
				res: httptest.NewRecorder(),
			}

			_, err := ctx.Body()
			require.Error(t, err)

			var errWithStatus fuego.ErrorWithStatus
			require.ErrorAs(t, err, &errWithStatus)
			require.Equal(t, http.StatusUnsupportedMediaType, errWithStatus.StatusCode())
		})
	}
}
