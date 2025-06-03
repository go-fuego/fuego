package fuego

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMockContext_MustBody(t *testing.T) {
	t.Run("can read body without error", func(t *testing.T) {
		ctx := NewMockContext[string, any]("test body", nil)
		body := ctx.MustBody()
		require.Equal(t, "test body", body)
	})
}

func TestMockContext_Params(t *testing.T) {
	t.Run("returns empty params by default", func(t *testing.T) {
		ctx := NewMockContextNoBody()
		params, err := ctx.Params()
		require.NoError(t, err)
		require.Nil(t, params)
	})
}

func TestMockContext_MustParams(t *testing.T) {
	t.Run("returns empty params by default", func(t *testing.T) {
		ctx := NewMockContextNoBody()
		params := ctx.MustParams()
		require.Nil(t, params)
	})
}

func TestMockContext_HasHeader(t *testing.T) {
	t.Run("returns false for non-existent header", func(t *testing.T) {
		ctx := NewMockContextNoBody()
		has := ctx.HasHeader("X-Test")
		require.False(t, has)
	})

	t.Run("returns true for existing header", func(t *testing.T) {
		ctx := NewMockContextNoBody()
		ctx.SetHeader("X-Test", "value")
		has := ctx.HasHeader("X-Test")
		require.True(t, has)
	})
}

func TestMockContext_HasCookie(t *testing.T) {
	t.Run("returns false for non-existent cookie", func(t *testing.T) {
		ctx := NewMockContextNoBody()
		has := ctx.HasCookie("session")
		require.False(t, has)
	})

	t.Run("returns true for existing cookie", func(t *testing.T) {
		ctx := NewMockContextNoBody()
		ctx.SetCookie(http.Cookie{Name: "session", Value: "abc123"})
		has := ctx.HasCookie("session")
		require.True(t, has)
	})
}

func TestMockContext_Header(t *testing.T) {
	t.Run("returns empty string for non-existent header", func(t *testing.T) {
		ctx := NewMockContextNoBody()
		header := ctx.Header("X-Test")
		require.Empty(t, header)
	})

	t.Run("returns header value for existing header", func(t *testing.T) {
		ctx := NewMockContextNoBody()
		ctx.SetHeader("X-Test", "test-value")
		header := ctx.Header("X-Test")
		require.Equal(t, "test-value", header)
	})
}

func TestMockContext_SetHeader(t *testing.T) {
	t.Run("sets header value", func(t *testing.T) {
		ctx := NewMockContextNoBody()
		ctx.SetHeader("X-Test", "test-value")
		header := ctx.Header("X-Test")
		require.Equal(t, "test-value", header)
	})
}

func TestMockContext_PathParam(t *testing.T) {
	t.Run("returns empty string for non-existent path param", func(t *testing.T) {
		ctx := NewMockContextNoBody()
		param := ctx.PathParam("id")
		require.Empty(t, param)
	})
}

func TestMockContext_PathParamIntErr(t *testing.T) {
	t.Run("returns error for non-existent path param", func(t *testing.T) {
		ctx := NewMockContextNoBody()
		_, err := ctx.PathParamIntErr("id")
		require.Error(t, err)
	})
}

func TestMockContext_PathParamInt(t *testing.T) {
	t.Run("returns 0 for non-existent path param", func(t *testing.T) {
		ctx := NewMockContextNoBody()
		param := ctx.PathParamInt("id")
		require.Zero(t, param)
	})
}

func TestMockContext_Request(t *testing.T) {
	t.Run("returns nil request by default", func(t *testing.T) {
		ctx := NewMockContextNoBody()
		req := ctx.Request()
		require.Nil(t, req)
	})
}

func TestMockContext_Response(t *testing.T) {
	t.Run("returns nil response by default", func(t *testing.T) {
		ctx := NewMockContextNoBody()
		res := ctx.Response()
		require.Nil(t, res)
	})
}

func TestMockContext_SetStatus(t *testing.T) {
	t.Run("sets status code without error", func(t *testing.T) {
		ctx := NewMockContextNoBody()
		// SetStatus should not panic when response is nil
		require.NotPanics(t, func() {
			ctx.SetStatus(201)
		})
	})
}

func TestMockContext_Cookie(t *testing.T) {
	t.Run("returns error for non-existent cookie", func(t *testing.T) {
		ctx := NewMockContextNoBody()
		_, err := ctx.Cookie("session")
		require.Error(t, err)
	})

	t.Run("returns cookie for existing cookie", func(t *testing.T) {
		ctx := NewMockContextNoBody()
		expectedCookie := http.Cookie{Name: "session", Value: "abc123"}
		ctx.SetCookie(expectedCookie)
		cookie, err := ctx.Cookie("session")
		require.NoError(t, err)
		require.Equal(t, expectedCookie.Name, cookie.Name)
		require.Equal(t, expectedCookie.Value, cookie.Value)
	})
}

func TestMockContext_SetCookie(t *testing.T) {
	t.Run("sets cookie", func(t *testing.T) {
		ctx := NewMockContextNoBody()
		cookie := http.Cookie{Name: "session", Value: "abc123"}
		ctx.SetCookie(cookie)

		retrievedCookie, err := ctx.Cookie("session")
		require.NoError(t, err)
		require.Equal(t, cookie.Name, retrievedCookie.Name)
		require.Equal(t, cookie.Value, retrievedCookie.Value)
	})
}

func TestMockContext_MainLang(t *testing.T) {
	t.Run("returns empty string by default", func(t *testing.T) {
		ctx := NewMockContextNoBody()
		lang := ctx.MainLang()
		require.Empty(t, lang)
	})
}

func TestMockContext_MainLocale(t *testing.T) {
	t.Run("returns empty string by default", func(t *testing.T) {
		ctx := NewMockContextNoBody()
		locale := ctx.MainLocale()
		require.Empty(t, locale)
	})
}

func TestMockContext_Redirect(t *testing.T) {
	t.Run("returns nil without error", func(t *testing.T) {
		ctx := NewMockContextNoBody()
		result, err := ctx.Redirect(301, "/test")
		require.NoError(t, err)
		require.Nil(t, result)
	})
}

func TestMockContext_Render(t *testing.T) {
	t.Run("panics as not implemented", func(t *testing.T) {
		ctx := NewMockContextNoBody()
		require.Panics(t, func() {
			ctx.Render("template.html", "data")
		})
	})
}

func TestMockContext_SetQueryParamBool(t *testing.T) {
	t.Run("sets boolean query parameter", func(t *testing.T) {
		ctx := NewMockContextNoBody()
		ctx.SetQueryParamBool("active", true)

		param := ctx.QueryParamBool("active")
		require.True(t, param)
	})

	t.Run("sets false boolean query parameter", func(t *testing.T) {
		ctx := NewMockContextNoBody()
		ctx.SetQueryParamBool("active", false)

		param := ctx.QueryParamBool("active")
		require.False(t, param)
	})
}

func TestMockContext_Body(t *testing.T) {
	t.Run("returns body", func(t *testing.T) {
		testBody := "test body"
		ctx := NewMockContext[string, any](testBody, nil)

		body, err := ctx.Body()
		require.NoError(t, err)
		require.Equal(t, testBody, body)
	})
}

func TestMockContext_SetQueryParam(t *testing.T) {
	t.Run("sets string query parameter", func(t *testing.T) {
		ctx := NewMockContextNoBody()
		ctx.SetQueryParam("name", "test")

		param := ctx.QueryParam("name")
		require.Equal(t, "test", param)
	})
}

func TestMockContext_SetQueryParamInt(t *testing.T) {
	t.Run("sets integer query parameter", func(t *testing.T) {
		ctx := NewMockContextNoBody()
		ctx.SetQueryParamInt("count", 42)

		param := ctx.QueryParamInt("count")
		require.Equal(t, 42, param)
	})
}

func TestMockContext_MainLang_WithHeader(t *testing.T) {
	t.Run("returns main language from Accept-Language header", func(t *testing.T) {
		ctx := NewMockContextNoBody()
		ctx.SetHeader("Accept-Language", "en-US,en;q=0.9,fr;q=0.8")

		lang := ctx.MainLang()
		require.Equal(t, "en", lang)
	})
}

func TestMockContext_MainLocale_WithHeader(t *testing.T) {
	t.Run("returns main locale from Accept-Language header", func(t *testing.T) {
		ctx := NewMockContextNoBody()
		ctx.SetHeader("Accept-Language", "en-US,en;q=0.9,fr;q=0.8")

		locale := ctx.MainLocale()
		require.Equal(t, "en-US,en;q=0.9,fr;q=0.8", locale)
	})
}

func TestMockContext_PathParamInt_WithValidParam(t *testing.T) {
	t.Run("returns valid integer path parameter", func(t *testing.T) {
		ctx := NewMockContextNoBody()
		ctx.PathParams["id"] = "123"

		param := ctx.PathParamInt("id")
		require.Equal(t, 123, param)
	})
}

func TestMockContext_PathParamIntErr_WithValidParam(t *testing.T) {
	t.Run("returns valid integer path parameter", func(t *testing.T) {
		ctx := NewMockContextNoBody()
		ctx.PathParams["id"] = "123"

		param, err := ctx.PathParamIntErr("id")
		require.NoError(t, err)
		require.Equal(t, 123, param)
	})

	t.Run("returns error for invalid integer", func(t *testing.T) {
		ctx := NewMockContextNoBody()
		ctx.PathParams["id"] = "invalid"

		_, err := ctx.PathParamIntErr("id")
		require.Error(t, err)
	})
}

func TestMockContext_PathParam_WithParam(t *testing.T) {
	t.Run("returns path parameter", func(t *testing.T) {
		ctx := NewMockContextNoBody()
		ctx.PathParams["name"] = "test"

		param := ctx.PathParam("name")
		require.Equal(t, "test", param)
	})
}
