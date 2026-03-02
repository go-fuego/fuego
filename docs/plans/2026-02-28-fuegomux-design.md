# fuegomux — gorilla/mux Adapter Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Create a gorilla/mux adapter for Fuego, enabling typed controllers, OpenAPI spec generation, and request validation with gorilla/mux routers.

**Architecture:** Follow the same adapter pattern as `extra/fuegogin` and `extra/fuegoecho` — a standalone Go module under `extra/fuegomux/` with a context wrapper, route registerer, and OpenAPI handler. gorilla/mux uses `{name}` path params (matching OpenAPI already) but supports `{name:regex}` — we strip regex for OpenAPI paths and extract patterns for parameter schemas.

**Tech Stack:** Go, gorilla/mux v1.8.1, fuego core (Engine, Registerer, Flow, CommonContext)

---

### Task 1: Scaffold fuegomux module

**Files:**
- Create: `extra/fuegomux/go.mod`

**Step 1: Create the module directory and go.mod**

```
extra/fuegomux/go.mod
```

```go
module github.com/go-fuego/fuego/extra/fuegomux

go 1.25.7

require (
	github.com/go-fuego/fuego v0.19.0
	github.com/gorilla/mux v1.8.1
	github.com/stretchr/testify v1.11.1
)
```

**Step 2: Add to go.work**

Add `./extra/fuegomux` to the `use` block in `go.work` (after `./extra/fuegogin`).

**Step 3: Run go mod tidy**

Run: `cd extra/fuegomux && go mod tidy`
Expected: go.sum created, dependencies resolved

**Step 4: Commit**

```bash
git add extra/fuegomux/go.mod extra/fuegomux/go.sum go.work
git commit -m "feat: scaffold fuegomux module for gorilla/mux adapter"
```

---

### Task 2: Implement context.go

**Files:**
- Create: `extra/fuegomux/context.go`

**Step 1: Write the failing test**

Create `extra/fuegomux/context_test.go`:

```go
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
```

**Step 2: Run test to verify it fails**

Run: `cd extra/fuegomux && go test ./... -v -run TestMuxContext`
Expected: FAIL — `muxContext` type not found

**Step 3: Write the context implementation**

Create `extra/fuegomux/context.go`:

```go
package fuegomux

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/gorilla/mux"

	"github.com/go-fuego/fuego"
	"github.com/go-fuego/fuego/internal"
)

type muxContext[B, P any] struct {
	internal.CommonContext[B]
	req *http.Request
	res http.ResponseWriter
}

var (
	_ fuego.Context[any, any]         = &muxContext[any, any]{}
	_ fuego.ContextWithBody[any]      = &muxContext[any, any]{}
	_ fuego.ContextFlowable[any, any] = &muxContext[any, any]{}
)

func (c *muxContext[B, P]) Body() (B, error) {
	var body B
	err := json.NewDecoder(c.req.Body).Decode(&body)
	if err != nil {
		return body, err
	}
	return fuego.TransformAndValidate(c, body)
}

func (c *muxContext[B, P]) MustBody() B {
	body, err := c.Body()
	if err != nil {
		panic(err)
	}
	return body
}

func (c *muxContext[B, P]) Params() (P, error) {
	var params P
	return params, nil
}

func (c *muxContext[B, P]) MustParams() P {
	params, err := c.Params()
	if err != nil {
		panic(err)
	}
	return params
}

func (c *muxContext[B, P]) Context() context.Context {
	return c.req.Context()
}

func (c *muxContext[B, P]) Cookie(name string) (*http.Cookie, error) {
	return c.req.Cookie(name)
}

func (c *muxContext[B, P]) HasCookie(name string) bool {
	_, err := c.Cookie(name)
	return err == nil
}

func (c *muxContext[B, P]) Header(key string) string {
	return c.req.Header.Get(key)
}

func (c *muxContext[B, P]) HasHeader(key string) bool {
	_, ok := c.req.Header[key]
	return ok
}

func (c *muxContext[B, P]) SetHeader(key, value string) {
	c.res.Header().Set(key, value)
}

func (c *muxContext[B, P]) SetCookie(cookie http.Cookie) {
	http.SetCookie(c.res, &cookie)
}

func (c *muxContext[B, P]) PathParam(name string) string {
	return mux.Vars(c.req)[name]
}

func (c *muxContext[B, P]) PathParamIntErr(name string) (int, error) {
	return fuego.PathParamIntErr(c, name)
}

func (c *muxContext[B, P]) PathParamInt(name string) int {
	param, _ := fuego.PathParamIntErr(c, name)
	return param
}

func (c *muxContext[B, P]) MainLang() string {
	return strings.Split(c.MainLocale(), "-")[0]
}

func (c *muxContext[B, P]) MainLocale() string {
	return strings.Split(c.req.Header.Get("Accept-Language"), ",")[0]
}

func (c *muxContext[B, P]) Redirect(code int, url string) (any, error) {
	http.Redirect(c.res, c.req, url, code)
	return nil, nil
}

func (c *muxContext[B, P]) Render(templateToExecute string, data any, templateGlobsToOverride ...string) (fuego.CtxRenderer, error) {
	panic("unimplemented")
}

func (c *muxContext[B, P]) Request() *http.Request {
	return c.req
}

func (c *muxContext[B, P]) Response() http.ResponseWriter {
	return c.res
}

func (c *muxContext[B, P]) SetStatus(code int) {
	c.res.WriteHeader(code)
}

func (c *muxContext[B, P]) Serialize(data any) error {
	return fuego.Send(c.res, c.req, data)
}

func (c *muxContext[B, P]) SerializeError(err error) {
	statusCode := http.StatusInternalServerError
	var errorWithStatusCode fuego.ErrorWithStatus
	if errors.As(err, &errorWithStatusCode) {
		statusCode = errorWithStatusCode.StatusCode()
	}
	c.res.Header().Set("Content-Type", "application/json")
	c.res.WriteHeader(statusCode)
	json.NewEncoder(c.res).Encode(err)
}

func (c *muxContext[B, P]) SetDefaultStatusCode() {
	if c.DefaultStatusCode == 0 {
		c.DefaultStatusCode = http.StatusOK
	}
	c.SetStatus(c.DefaultStatusCode)
}
```

**Step 4: Run tests to verify they pass**

Run: `cd extra/fuegomux && go test ./... -v -run TestMuxContext`
Expected: PASS

**Step 5: Commit**

```bash
git add extra/fuegomux/context.go extra/fuegomux/context_test.go
git commit -m "feat(fuegomux): implement muxContext with gorilla/mux path param support"
```

---

### Task 3: Implement adaptor.go — path conversion and route registration

**Files:**
- Create: `extra/fuegomux/adaptor.go`
- Create: `extra/fuegomux/adaptor_test.go`

**Step 1: Write the failing test for path conversion**

Add to `extra/fuegomux/adaptor_test.go`:

```go
package fuegomux

import (
	"testing"

	"github.com/stretchr/testify/assert"
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
```

**Step 2: Run test to verify it fails**

Run: `cd extra/fuegomux && go test ./... -v -run "TestMuxToFuegoRoute|TestExtractPathParamPatterns"`
Expected: FAIL — functions not found

**Step 3: Write path conversion + full adaptor.go**

Create `extra/fuegomux/adaptor.go`:

```go
package fuegomux

import (
	"net/http"
	"regexp"

	"github.com/gorilla/mux"

	"github.com/go-fuego/fuego"
	"github.com/go-fuego/fuego/internal"
)

// pathRegex matches gorilla/mux path params with regex constraints: {name:pattern}
var pathRegex = regexp.MustCompile(`\{([a-zA-Z_][a-zA-Z0-9_]*):([^}]+)\}`)

// muxToFuegoRoute strips regex constraints from gorilla/mux paths.
// {id:[0-9]+} → {id}
func muxToFuegoRoute(path string) string {
	return pathRegex.ReplaceAllString(path, `{$1}`)
}

// extractPathParamPatterns extracts regex patterns from path params.
// Returns a map of param name → regex pattern for params that have constraints.
func extractPathParamPatterns(path string) map[string]string {
	patterns := make(map[string]string)
	matches := pathRegex.FindAllStringSubmatch(path, -1)
	for _, m := range matches {
		patterns[m[1]] = m[2]
	}
	return patterns
}

// OpenAPIHandler implements fuego.OpenAPIServable for gorilla/mux routers.
type OpenAPIHandler struct {
	Router *mux.Router
}

func (o *OpenAPIHandler) SpecHandler(e *fuego.Engine) {
	Get(e, o.Router, e.OpenAPI.Config.SpecURL, e.SpecHandler(), fuego.OptionHide(), fuego.OptionMiddleware(e.OpenAPI.Config.SwaggerMiddlewares...))
}

func (o *OpenAPIHandler) UIHandler(e *fuego.Engine) {
	GetMux(
		e,
		o.Router,
		e.OpenAPI.Config.SwaggerURL,
		e.OpenAPI.Config.UIHandler(e.OpenAPI.Config.SpecURL).ServeHTTP,
		fuego.OptionHide(),
		fuego.OptionMiddleware(e.OpenAPI.Config.SwaggerMiddlewares...),
	)
}

// --- Native mux handler registration (Level 1 & 2) ---

func GetMux(engine *fuego.Engine, muxRouter *mux.Router, path string, handler http.HandlerFunc, options ...func(*fuego.BaseRoute)) *fuego.Route[any, any, any] {
	return handleMux(engine, muxRouter, http.MethodGet, path, handler, options...)
}

func PostMux(engine *fuego.Engine, muxRouter *mux.Router, path string, handler http.HandlerFunc, options ...func(*fuego.BaseRoute)) *fuego.Route[any, any, any] {
	return handleMux(engine, muxRouter, http.MethodPost, path, handler, options...)
}

func PutMux(engine *fuego.Engine, muxRouter *mux.Router, path string, handler http.HandlerFunc, options ...func(*fuego.BaseRoute)) *fuego.Route[any, any, any] {
	return handleMux(engine, muxRouter, http.MethodPut, path, handler, options...)
}

func DeleteMux(engine *fuego.Engine, muxRouter *mux.Router, path string, handler http.HandlerFunc, options ...func(*fuego.BaseRoute)) *fuego.Route[any, any, any] {
	return handleMux(engine, muxRouter, http.MethodDelete, path, handler, options...)
}

func PatchMux(engine *fuego.Engine, muxRouter *mux.Router, path string, handler http.HandlerFunc, options ...func(*fuego.BaseRoute)) *fuego.Route[any, any, any] {
	return handleMux(engine, muxRouter, http.MethodPatch, path, handler, options...)
}

func OptionsMux(engine *fuego.Engine, muxRouter *mux.Router, path string, handler http.HandlerFunc, options ...func(*fuego.BaseRoute)) *fuego.Route[any, any, any] {
	return handleMux(engine, muxRouter, http.MethodOptions, path, handler, options...)
}

// --- Fuego typed handler registration (Level 3 & 4) ---

func Get[T, B, P any](engine *fuego.Engine, muxRouter *mux.Router, path string, handler func(c fuego.Context[B, P]) (T, error), options ...func(*fuego.BaseRoute)) *fuego.Route[T, B, P] {
	return handleFuego(engine, muxRouter, http.MethodGet, path, handler, options...)
}

func Post[T, B, P any](engine *fuego.Engine, muxRouter *mux.Router, path string, handler func(c fuego.Context[B, P]) (T, error), options ...func(*fuego.BaseRoute)) *fuego.Route[T, B, P] {
	return handleFuego(engine, muxRouter, http.MethodPost, path, handler, options...)
}

func Put[T, B, P any](engine *fuego.Engine, muxRouter *mux.Router, path string, handler func(c fuego.Context[B, P]) (T, error), options ...func(*fuego.BaseRoute)) *fuego.Route[T, B, P] {
	return handleFuego(engine, muxRouter, http.MethodPut, path, handler, options...)
}

func Delete[T, B, P any](engine *fuego.Engine, muxRouter *mux.Router, path string, handler func(c fuego.Context[B, P]) (T, error), options ...func(*fuego.BaseRoute)) *fuego.Route[T, B, P] {
	return handleFuego(engine, muxRouter, http.MethodDelete, path, handler, options...)
}

func Patch[T, B, P any](engine *fuego.Engine, muxRouter *mux.Router, path string, handler func(c fuego.Context[B, P]) (T, error), options ...func(*fuego.BaseRoute)) *fuego.Route[T, B, P] {
	return handleFuego(engine, muxRouter, http.MethodPatch, path, handler, options...)
}

func Options[T, B, P any](engine *fuego.Engine, muxRouter *mux.Router, path string, handler func(c fuego.Context[B, P]) (T, error), options ...func(*fuego.BaseRoute)) *fuego.Route[T, B, P] {
	return handleFuego(engine, muxRouter, http.MethodOptions, path, handler, options...)
}

// --- Internal registration ---

func handleFuego[T, B, P any](engine *fuego.Engine, muxRouter *mux.Router, method, path string, fuegoHandler func(c fuego.Context[B, P]) (T, error), options ...func(*fuego.BaseRoute)) *fuego.Route[T, B, P] {
	baseRoute := fuego.NewBaseRoute(method, muxToFuegoRoute(path), fuegoHandler, engine, options...)
	return fuego.Registers(engine, muxRouteRegisterer[T, B, P]{
		muxRouter:    muxRouter,
		route:        fuego.Route[T, B, P]{BaseRoute: baseRoute},
		httpHandler:  MuxHandler(engine, fuegoHandler, baseRoute),
		originalPath: path,
	})
}

func handleMux(engine *fuego.Engine, muxRouter *mux.Router, method, path string, handler http.HandlerFunc, options ...func(*fuego.BaseRoute)) *fuego.Route[any, any, any] {
	baseRoute := fuego.NewBaseRoute(method, muxToFuegoRoute(path), handler, engine, options...)
	return fuego.Registers(engine, muxRouteRegisterer[any, any, any]{
		muxRouter:    muxRouter,
		route:        fuego.Route[any, any, any]{BaseRoute: baseRoute},
		httpHandler:  handler,
		originalPath: path,
	})
}

// --- Route registerer ---

type muxRouteRegisterer[T, B, P any] struct {
	muxRouter    *mux.Router
	httpHandler  http.HandlerFunc
	route        fuego.Route[T, B, P]
	originalPath string
}

func (a muxRouteRegisterer[T, B, P]) Register() fuego.Route[T, B, P] {
	muxRoute := a.muxRouter.HandleFunc(a.originalPath, a.httpHandler).Methods(a.route.Method)

	// Get the full path template including any subrouter prefix
	if tpl, err := muxRoute.GetPathTemplate(); err == nil {
		a.route.Path = muxToFuegoRoute(tpl)
	}

	return a.route
}

// MuxHandler converts a Fuego handler to an http.HandlerFunc.
func MuxHandler[B, T, P any](engine *fuego.Engine, handler func(c fuego.Context[B, P]) (T, error), route fuego.BaseRoute) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := &muxContext[B, P]{
			CommonContext: internal.CommonContext[B]{
				CommonCtx:         r.Context(),
				UrlValues:         r.URL.Query(),
				OpenAPIParams:     route.Params,
				DefaultStatusCode: route.DefaultStatusCode,
			},
			req: r,
			res: w,
		}
		fuego.Flow(engine, ctx, handler)
	}
}
```

**Step 4: Write tests for route registration and OpenAPI spec**

Add to `extra/fuegomux/adaptor_test.go`:

```go
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
```

**Step 5: Run tests**

Run: `cd extra/fuegomux && go test ./... -v`
Expected: ALL PASS

**Step 6: Commit**

```bash
git add extra/fuegomux/adaptor.go extra/fuegomux/adaptor_test.go
git commit -m "feat(fuegomux): implement route registration with regex path support"
```

---

### Task 4: Integration test — full HTTP round-trip

**Files:**
- Modify: `extra/fuegomux/adaptor_test.go`

**Step 1: Write integration test**

Add to `extra/fuegomux/adaptor_test.go`:

```go
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
```

Add `"net/http/httptest"` and `"strings"` to imports.

**Step 2: Run test**

Run: `cd extra/fuegomux && go test ./... -v -run "TestFuegoHandler"`
Expected: PASS

**Step 3: Commit**

```bash
git add extra/fuegomux/adaptor_test.go
git commit -m "test(fuegomux): add HTTP round-trip integration tests"
```

---

### Task 5: Create example app

**Files:**
- Create: `examples/mux-compat/go.mod`
- Create: `examples/mux-compat/main.go`
- Create: `examples/mux-compat/handlers.go`
- Create: `examples/mux-compat/main_test.go`

**Step 1: Create go.mod**

```go
module github.com/go-fuego/fuego/examples/mux-compat

go 1.25.7

replace github.com/go-fuego/fuego => ../..

replace github.com/go-fuego/fuego/extra/fuegomux => ../../extra/fuegomux

require (
	github.com/go-fuego/fuego v0.19.0
	github.com/go-fuego/fuego/extra/fuegomux v0.0.0
	github.com/gorilla/mux v1.8.1
	github.com/stretchr/testify v1.11.1
)
```

**Step 2: Create handlers.go**

```go
package main

import (
	"fmt"
	"net/http"

	"github.com/go-fuego/fuego"
)

func muxController(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "pong")
}

func fuegoControllerGet(c fuego.ContextNoBody) (HelloResponse, error) {
	return HelloResponse{
		Message: "Hello",
	}, nil
}

func fuegoControllerPost(c fuego.ContextWithBody[HelloRequest]) (*HelloResponse, error) {
	body, err := c.Body()
	if err != nil {
		return nil, err
	}

	if body.Word == "forbidden" {
		return nil, fuego.BadRequestError{Title: "Forbidden word"}
	}

	name := c.QueryParam("name")

	return &HelloResponse{
		Message: fmt.Sprintf("Hello %s, %s", body.Word, name),
	}, nil
}
```

**Step 3: Create main.go**

```go
package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/gorilla/mux"

	"github.com/go-fuego/fuego"
	"github.com/go-fuego/fuego/extra/fuegomux"
	"github.com/go-fuego/fuego/option"
	"github.com/go-fuego/fuego/param"
)

type HelloRequest struct {
	Word string `json:"word" validate:"required,min=2"`
}

var _ fuego.InTransformer = &HelloRequest{}

type HelloResponse struct {
	Message string `json:"message"`
}

func main() {
	r, _ := server()

	fmt.Println("OpenAPI at http://localhost:8980/swagger")

	err := http.ListenAndServe(":8980", r)
	if err != nil {
		panic(err)
	}
}

func server() (*mux.Router, *fuego.OpenAPI) {
	muxRouter := mux.NewRouter()
	engine := fuego.NewEngine()

	// Register native mux controller
	muxRouter.HandleFunc("/mux", muxController).Methods(http.MethodGet)

	// 1. Level 1: Register native mux controller with OpenAPI spec
	fuegomux.GetMux(engine, muxRouter, "/mux-with-openapi", muxController)

	// 2. Level 2: Native mux controller with OpenAPI options
	fuegomux.GetMux(engine, muxRouter, "/mux-with-openapi-and-options", muxController,
		option.Summary("Mux controller with options"),
		option.Description("Some description"),
		option.OperationID("MyCustomOperationID"),
		option.Tags("Mux"),
	)

	// 3. Level 3: Fuego controller with gorilla/mux router
	fuegomux.Get(engine, muxRouter, "/fuego", fuegoControllerGet)

	// 4. Level 4: Fuego controller with options
	fuegomux.Post(engine, muxRouter, "/fuego-with-options", fuegoControllerPost,
		option.Description("Some description"),
		option.OperationID("SomeOperationID"),
		option.AddError(409, "Name Already Exists"),
		option.DefaultStatusCode(201),
		option.Tags("Fuego"),
		option.Query("name", "Your name", param.Example("name example", "John Carmack")),
		option.Header("X-Request-ID", "Request ID", param.Default("123456")),
		option.Header("Content-Type", "Content Type", param.Default("application/json")),
	)

	// Groups & path parameters with regex
	sub := muxRouter.PathPrefix("/my-group/{id:[0-9]+}").Subrouter()
	fuegomux.Get(engine, sub, "/fuego", fuegoControllerGet,
		option.Summary("Route with subrouter and id"),
		option.Tags("Fuego"),
	)

	engine.RegisterOpenAPIRoutes(&fuegomux.OpenAPIHandler{Router: muxRouter})

	return muxRouter, engine.OpenAPI
}

func (h *HelloRequest) InTransform(ctx context.Context) error {
	h.Word = strings.ToLower(h.Word)

	if h.Word == "apple" {
		return fuego.BadRequestError{Title: "Word not allowed", Err: errors.New("forbidden word"), Detail: "The word 'apple' is not allowed"}
	}

	if h.Word == "banana" {
		return errors.New("banana is not allowed")
	}

	if user := ctx.Value("user"); user == "secret agent" {
		h.Word = "*****"
	}

	return nil
}
```

**Step 4: Create main_test.go**

```go
package main

import (
	"testing"

	"github.com/go-fuego/fuego"
	"github.com/stretchr/testify/require"
)

func TestFuegoControllerPost(t *testing.T) {
	testCtx := fuego.NewMockContext(HelloRequest{Word: "World"}, any(nil))
	testCtx.QueryParams().Set("name", "Ewen")

	response, err := fuegoControllerPost(testCtx)
	require.NoError(t, err)
	require.Equal(t, "Hello World, Ewen", response.Message)
}
```

**Step 5: Add to go.work**

Add `./examples/mux-compat` to the `use` block in `go.work`.

**Step 6: Run go mod tidy and tests**

Run: `cd examples/mux-compat && go mod tidy && go test ./... -v`
Expected: PASS

**Step 7: Commit**

```bash
git add examples/mux-compat/ go.work
git commit -m "feat(fuegomux): add mux-compat example app"
```

---

### Task 6: Final verification — run all tests across workspace

**Step 1: Run fuegomux tests**

Run: `cd extra/fuegomux && go test ./... -v`
Expected: ALL PASS

**Step 2: Run example tests**

Run: `cd examples/mux-compat && go test ./... -v`
Expected: ALL PASS

**Step 3: Run linter on new code**

Run: `cd extra/fuegomux && goimports -w . && go vet ./...`
Expected: No issues

**Step 4: Verify existing tests are not broken**

Run: `cd ~/dev/fuego && go test ./... -v -count=1`
Expected: ALL PASS (existing tests unaffected)

**Step 5: Final commit if any formatting changes**

```bash
git add -A && git diff --cached --quiet || git commit -m "style: format fuegomux code"
```
