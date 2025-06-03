package fuegogin

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/go-fuego/fuego"
	"github.com/go-fuego/fuego/internal"
)

func setupGinContext() (*gin.Context, *httptest.ResponseRecorder) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	return c, w
}

func TestGinContext_Body(t *testing.T) {
	t.Run("can read JSON body", func(t *testing.T) {
		ginCtx, _ := setupGinContext()
		ginCtx.Request = httptest.NewRequest("POST", "/", strings.NewReader(`{"name":"test","value":123}`))
		ginCtx.Request.Header.Set("Content-Type", "application/json")

		ctx := ginContext[map[string]any, any]{
			CommonContext: internal.CommonContext[map[string]any]{},
			ginCtx:        ginCtx,
		}

		body, err := ctx.Body()
		require.NoError(t, err)
		require.Equal(t, "test", body["name"])
		require.Equal(t, float64(123), body["value"]) // JSON numbers are float64
	})

	t.Run("returns error for invalid JSON", func(t *testing.T) {
		ginCtx, _ := setupGinContext()
		ginCtx.Request = httptest.NewRequest("POST", "/", strings.NewReader(`invalid json`))
		ginCtx.Request.Header.Set("Content-Type", "application/json")

		ctx := ginContext[map[string]any, any]{
			CommonContext: internal.CommonContext[map[string]any]{},
			ginCtx:        ginCtx,
		}

		_, err := ctx.Body()
		require.Error(t, err)
	})
}

func TestGinContext_Context(t *testing.T) {
	ginCtx, _ := setupGinContext()
	ctx := ginContext[any, any]{
		ginCtx: ginCtx,
	}

	result := ctx.Context()
	require.Equal(t, ginCtx, result)
}

func TestGinContext_Cookie(t *testing.T) {
	t.Run("existing cookie", func(t *testing.T) {
		ginCtx, _ := setupGinContext()
		ginCtx.Request = httptest.NewRequest("GET", "/", nil)
		ginCtx.Request.AddCookie(&http.Cookie{Name: "test", Value: "value"})

		ctx := ginContext[any, any]{ginCtx: ginCtx}

		cookie, err := ctx.Cookie("test")
		require.NoError(t, err)
		require.Equal(t, "test", cookie.Name)
		require.Equal(t, "value", cookie.Value)
	})

	t.Run("non-existing cookie", func(t *testing.T) {
		ginCtx, _ := setupGinContext()
		ginCtx.Request = httptest.NewRequest("GET", "/", nil)

		ctx := ginContext[any, any]{ginCtx: ginCtx}

		_, err := ctx.Cookie("nonexistent")
		require.Error(t, err)
	})
}

func TestGinContext_Header(t *testing.T) {
	ginCtx, _ := setupGinContext()
	ginCtx.Request = httptest.NewRequest("GET", "/", nil)
	ginCtx.Request.Header.Set("X-Test", "test-value")

	ctx := ginContext[any, any]{ginCtx: ginCtx}

	result := ctx.Header("X-Test")
	require.Equal(t, "test-value", result)

	result = ctx.Header("Non-Existent")
	require.Empty(t, result)
}

func TestGinContext_MustBody(t *testing.T) {
	t.Run("successful body read", func(t *testing.T) {
		ginCtx, _ := setupGinContext()
		ginCtx.Request = httptest.NewRequest("POST", "/", strings.NewReader(`{"name":"test"}`))
		ginCtx.Request.Header.Set("Content-Type", "application/json")

		ctx := ginContext[map[string]any, any]{
			CommonContext: internal.CommonContext[map[string]any]{},
			ginCtx:        ginCtx,
		}

		body := ctx.MustBody()
		require.Equal(t, "test", body["name"])
	})

	t.Run("panics on error", func(t *testing.T) {
		ginCtx, _ := setupGinContext()
		ginCtx.Request = httptest.NewRequest("POST", "/", strings.NewReader(`invalid json`))
		ginCtx.Request.Header.Set("Content-Type", "application/json")

		ctx := ginContext[map[string]any, any]{
			CommonContext: internal.CommonContext[map[string]any]{},
			ginCtx:        ginCtx,
		}

		require.Panics(t, func() {
			ctx.MustBody()
		})
	})
}

func TestGinContext_Params(t *testing.T) {
	ginCtx, _ := setupGinContext()
	ctx := ginContext[any, struct{ ID int }]{ginCtx: ginCtx}

	params, err := ctx.Params()
	require.NoError(t, err)
	require.Zero(t, params.ID) // Default zero value
}

func TestGinContext_MustParams(t *testing.T) {
	ginCtx, _ := setupGinContext()
	ctx := ginContext[any, struct{ ID int }]{ginCtx: ginCtx}

	params := ctx.MustParams()
	require.Zero(t, params.ID) // Default zero value
}

func TestGinContext_PathParam(t *testing.T) {
	ginCtx, _ := setupGinContext()
	ginCtx.Params = gin.Params{
		{Key: "id", Value: "123"},
		{Key: "name", Value: "test"},
	}

	ctx := ginContext[any, any]{ginCtx: ginCtx}

	result := ctx.PathParam("id")
	require.Equal(t, "123", result)

	result = ctx.PathParam("nonexistent")
	require.Empty(t, result)
}

func TestGinContext_PathParamIntErr(t *testing.T) {
	t.Run("valid integer", func(t *testing.T) {
		ginCtx, _ := setupGinContext()
		ginCtx.Params = gin.Params{{Key: "id", Value: "123"}}

		ctx := ginContext[any, any]{ginCtx: ginCtx}

		result, err := ctx.PathParamIntErr("id")
		require.NoError(t, err)
		require.Equal(t, 123, result)
	})

	t.Run("invalid integer", func(t *testing.T) {
		ginCtx, _ := setupGinContext()
		ginCtx.Params = gin.Params{{Key: "id", Value: "invalid"}}

		ctx := ginContext[any, any]{ginCtx: ginCtx}

		_, err := ctx.PathParamIntErr("id")
		require.Error(t, err)
	})
}

func TestGinContext_PathParamInt(t *testing.T) {
	t.Run("valid integer", func(t *testing.T) {
		ginCtx, _ := setupGinContext()
		ginCtx.Params = gin.Params{{Key: "id", Value: "123"}}

		ctx := ginContext[any, any]{ginCtx: ginCtx}

		result := ctx.PathParamInt("id")
		require.Equal(t, 123, result)
	})

	t.Run("invalid integer returns 0", func(t *testing.T) {
		ginCtx, _ := setupGinContext()
		ginCtx.Params = gin.Params{{Key: "id", Value: "invalid"}}

		ctx := ginContext[any, any]{ginCtx: ginCtx}

		result := ctx.PathParamInt("id")
		require.Zero(t, result)
	})
}

func TestGinContext_MainLang(t *testing.T) {
	ginCtx, _ := setupGinContext()
	ginCtx.Request = httptest.NewRequest("GET", "/", nil)
	ginCtx.Request.Header.Set("Accept-Language", "en-US,en;q=0.9,fr;q=0.8")

	ctx := ginContext[any, any]{ginCtx: ginCtx}

	result := ctx.MainLang()
	require.Equal(t, "en", result)
}

func TestGinContext_MainLocale(t *testing.T) {
	ginCtx, _ := setupGinContext()
	ginCtx.Request = httptest.NewRequest("GET", "/", nil)
	ginCtx.Request.Header.Set("Accept-Language", "en-US,en;q=0.9,fr;q=0.8")

	ctx := ginContext[any, any]{ginCtx: ginCtx}

	result := ctx.MainLocale()
	require.Equal(t, "en-US", result)
}

func TestGinContext_Redirect(t *testing.T) {
	ginCtx, w := setupGinContext()
	ginCtx.Request = httptest.NewRequest("GET", "/", nil)
	ctx := ginContext[any, any]{ginCtx: ginCtx}

	result, err := ctx.Redirect(http.StatusFound, "/redirect")
	require.NoError(t, err)
	require.Nil(t, result)
	require.Equal(t, http.StatusFound, w.Code)
	require.Equal(t, "/redirect", w.Header().Get("Location"))
}

func TestGinContext_Render(t *testing.T) {
	ginCtx, _ := setupGinContext()
	ctx := ginContext[any, any]{ginCtx: ginCtx}

	require.Panics(t, func() {
		ctx.Render("template", "data")
	})
}

func TestGinContext_Request(t *testing.T) {
	ginCtx, _ := setupGinContext()
	req := httptest.NewRequest("GET", "/", nil)
	ginCtx.Request = req

	ctx := ginContext[any, any]{ginCtx: ginCtx}

	result := ctx.Request()
	require.Equal(t, req, result)
}

func TestGinContext_Response(t *testing.T) {
	ginCtx, _ := setupGinContext()
	ctx := ginContext[any, any]{ginCtx: ginCtx}

	result := ctx.Response()
	require.NotNil(t, result)
	require.Equal(t, ginCtx.Writer, result)
}

func TestGinContext_SetCookie(t *testing.T) {
	ginCtx, w := setupGinContext()
	ctx := ginContext[any, any]{ginCtx: ginCtx}

	cookie := http.Cookie{
		Name:     "test",
		Value:    "value",
		MaxAge:   3600,
		Path:     "/",
		Domain:   "example.com",
		Secure:   true,
		HttpOnly: true,
	}

	ctx.SetCookie(cookie)

	// Check if cookie was set in response
	cookies := w.Header()["Set-Cookie"]
	require.NotEmpty(t, cookies)
	require.Contains(t, cookies[0], "test=value")
}

func TestGinContext_HasCookie(t *testing.T) {
	t.Run("existing cookie", func(t *testing.T) {
		ginCtx, _ := setupGinContext()
		ginCtx.Request = httptest.NewRequest("GET", "/", nil)
		ginCtx.Request.AddCookie(&http.Cookie{Name: "test", Value: "value"})

		ctx := ginContext[any, any]{ginCtx: ginCtx}

		result := ctx.HasCookie("test")
		require.True(t, result)
	})

	t.Run("non-existing cookie", func(t *testing.T) {
		ginCtx, _ := setupGinContext()
		ginCtx.Request = httptest.NewRequest("GET", "/", nil)

		ctx := ginContext[any, any]{ginCtx: ginCtx}

		result := ctx.HasCookie("nonexistent")
		require.False(t, result)
	})
}

func TestGinContext_HasHeader(t *testing.T) {
	t.Run("existing header", func(t *testing.T) {
		ginCtx, _ := setupGinContext()
		ginCtx.Request = httptest.NewRequest("GET", "/", nil)
		ginCtx.Request.Header.Set("X-Test", "value")

		ctx := ginContext[any, any]{ginCtx: ginCtx}

		result := ctx.HasHeader("X-Test")
		require.True(t, result)
	})

	t.Run("non-existing header", func(t *testing.T) {
		ginCtx, _ := setupGinContext()
		ginCtx.Request = httptest.NewRequest("GET", "/", nil)

		ctx := ginContext[any, any]{ginCtx: ginCtx}

		result := ctx.HasHeader("X-Nonexistent")
		require.False(t, result)
	})
}

func TestGinContext_SetHeader(t *testing.T) {
	ginCtx, w := setupGinContext()
	ctx := ginContext[any, any]{ginCtx: ginCtx}

	ctx.SetHeader("X-Test", "test-value")

	result := w.Header().Get("X-Test")
	require.Equal(t, "test-value", result)
}

func TestGinContext_SetStatus(t *testing.T) {
	ginCtx, _ := setupGinContext()
	ctx := ginContext[any, any]{ginCtx: ginCtx}

	ctx.SetStatus(http.StatusCreated)

	// Gin doesn't write the status until a response is written
	require.Equal(t, http.StatusCreated, ginCtx.Writer.Status())
}

func TestGinContext_Serialize(t *testing.T) {
	ginCtx, w := setupGinContext()
	ctx := ginContext[any, any]{ginCtx: ginCtx}

	data := map[string]string{"message": "hello"}
	err := ctx.Serialize(data)

	require.NoError(t, err)
	require.Equal(t, http.StatusOK, w.Code)
	require.Contains(t, w.Body.String(), `"message":"hello"`)
}

func TestGinContext_SerializeError(t *testing.T) {
	t.Run("regular error", func(t *testing.T) {
		ginCtx, w := setupGinContext()
		ctx := ginContext[any, any]{ginCtx: ginCtx}

		err := assert.AnError
		ctx.SerializeError(err)

		require.Equal(t, http.StatusInternalServerError, w.Code)
		// Gin serializes errors as JSON objects, but the content may be empty for some errors
		require.NotEmpty(t, w.Body.String())
	})

	t.Run("error with status code", func(t *testing.T) {
		ginCtx, w := setupGinContext()
		ctx := ginContext[any, any]{ginCtx: ginCtx}

		err := fuego.BadRequestError{Err: assert.AnError}
		ctx.SerializeError(err)

		require.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestGinContext_SetDefaultStatusCode(t *testing.T) {
	t.Run("with default status code", func(t *testing.T) {
		ginCtx, _ := setupGinContext()
		ctx := ginContext[any, any]{
			CommonContext: internal.CommonContext[any]{
				DefaultStatusCode: http.StatusCreated,
			},
			ginCtx: ginCtx,
		}

		ctx.SetDefaultStatusCode()

		require.Equal(t, http.StatusCreated, ginCtx.Writer.Status())
	})

	t.Run("without default status code", func(t *testing.T) {
		ginCtx, _ := setupGinContext()
		ctx := ginContext[any, any]{
			CommonContext: internal.CommonContext[any]{
				DefaultStatusCode: 0,
			},
			ginCtx: ginCtx,
		}

		// Should not panic
		require.NotPanics(t, func() {
			ctx.SetDefaultStatusCode()
		})

		// When default status code is 0, it gets set to 200 and then applied
		require.Equal(t, http.StatusOK, ginCtx.Writer.Status())
	})
}
