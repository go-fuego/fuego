package fuegogin

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/go-fuego/fuego"
	"github.com/go-fuego/fuego/option"
	"github.com/stretchr/testify/require"
	"gotest.tools/v3/assert"
)

func orderMiddleware(s string) func(*gin.Context) {
	return func(c *gin.Context) {
		c.Request.Header.Add("X-Test-Order", s)
		c.Next()
	}
}

func TestGinPathWithGinGroup(t *testing.T) {
	basePath := "/api"
	apiPath := "/path"
	completePath := basePath + apiPath

	e := fuego.NewEngine()
	ginRouter := gin.New()
	group := ginRouter.Group(basePath)

	GetGin(e, group, apiPath, gin.HandlerFunc(func(ctx *gin.Context) {
		result, err := func(ctx *gin.Context) (string, error) { return "ok", nil }(ctx)
		if err != nil {
			ctx.JSON(500, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(200, gin.H{"result": result})
	}))

	routes := ginRouter.Routes()
	for _, route := range routes {
		assert.Equal(t, route.Path, completePath)
	}
}

func TestFuegoPathWithGinPathParam(t *testing.T) {
	basePath := "/api"
	apiPath := "/path/:id"
	e := fuego.NewEngine()
	ginRouter := gin.New()
	group := ginRouter.Group(basePath)

	Get(e, group, apiPath, func(c fuego.ContextWithBody[any]) (string, error) { return "ok", nil })

	spec := e.OutputOpenAPISpec()
	specPath := spec.Paths.Find("/api/path/{id}")

	routes := ginRouter.Routes()
	assert.Assert(t, len(routes) == 1, "Expected exactly one route registered in Gin")
	assert.Equal(t, routes[0].Path, basePath+apiPath)

	assert.Check(t, specPath != nil,
		"Expected path '/api/path/{id}' to be registered in OpenAPI spec")
}

func TestFuegoPathWithGinGroup(t *testing.T) {
	basePath := "/api"
	apiPath := "/path"
	completePath := basePath + apiPath

	e := fuego.NewEngine()
	ginRouter := gin.New()
	group := ginRouter.Group(basePath)

	Get(e, group, apiPath, func(c fuego.ContextWithBody[any]) (string, error) { return "ok", nil })

	routes := ginRouter.Routes()
	for _, route := range routes {
		assert.Equal(t, route.Path, completePath)
	}
}

func TestGinMiddlewares(t *testing.T) {
	t.Run("gin handler & one middleware", func(t *testing.T) {
		engine := fuego.NewEngine()
		router := gin.Default()

		GetGin(engine, router, "/test",
			func(ctx *gin.Context) {
				ctx.AbortWithStatus(200)
			},
			option.Middleware(orderMiddleware("First!")),
		)

		r := httptest.NewRequest(http.MethodGet, "/test", nil)
		r.Header.Set("X-Test-Order", "Start!")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, r)

		require.Equal(t, []string{"Start!", "First!"}, r.Header["X-Test-Order"])
	})

	t.Run("gin handler & multiple middlewares", func(t *testing.T) {
		engine := fuego.NewEngine()
		router := gin.Default()

		GetGin(engine, router, "/test",
			func(ctx *gin.Context) {
				ctx.AbortWithStatus(200)
			},
			option.Middleware(orderMiddleware("First!")),
			option.Middleware(orderMiddleware("Second!")),
		)

		r := httptest.NewRequest(http.MethodGet, "/test", nil)
		r.Header.Set("X-Test-Order", "Start!")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, r)

		require.Equal(t, []string{"Start!", "First!", "Second!"}, r.Header["X-Test-Order"])
	})

	t.Run("gin handler & group middlewares", func(t *testing.T) {
		engine := fuego.NewEngine()
		router := gin.Default()

		router.Use(orderMiddleware("First!"))

		GetGin(engine, router, "/test",
			func(ctx *gin.Context) {
				ctx.AbortWithStatus(200)
			},
			option.Middleware(orderMiddleware("Second!")),
			option.Middleware(orderMiddleware("Third!")),
		)

		r := httptest.NewRequest(http.MethodGet, "/test", nil)
		r.Header.Set("X-Test-Order", "Start!")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, r)

		require.Equal(t, []string{"Start!", "First!", "Second!", "Third!"}, r.Header["X-Test-Order"])
	})

	t.Run("fuego handler & multiple middlewares", func(t *testing.T) {
		engine := fuego.NewEngine()
		router := gin.Default()

		Get(engine, router, "/test",
			func(ctx fuego.ContextNoBody) (any, error) {
				return nil, nil
			},
			option.Middleware(orderMiddleware("First!")),
			option.Middleware(orderMiddleware("Second!")),
		)

		r := httptest.NewRequest(http.MethodGet, "/test", nil)
		r.Header.Set("X-Test-Order", "Start!")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, r)

		require.Equal(t, []string{"Start!", "First!", "Second!"}, r.Header["X-Test-Order"])
	})

	t.Run("panics on wrong middleware", func(t *testing.T) {
		engine := fuego.NewEngine()
		router := gin.Default()

		require.Panics(t, func() {
			GetGin(engine, router, "/test",
				func(ctx *gin.Context) {
					ctx.AbortWithStatus(200)
				},
				option.Middleware(func(next http.Handler) http.Handler {
					return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
						next.ServeHTTP(w, r)
					})
				}),
			)
		})
	})
}
