package fuego

import (
	"context"
	"errors"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestContext_PathParam(t *testing.T) {
	t.Run("can read path param", func(t *testing.T) {
		t.Skip("TODO: coming in go1.22")

		s := NewServer()
		Get(s, "/foo/{id}", func(c ContextNoBody) (ans, error) {
			return ans{Ans: c.PathParam("id")}, nil
		})

		r := httptest.NewRequest("GET", "/foo/123", nil)
		w := httptest.NewRecorder()

		s.Mux.ServeHTTP(w, r)

		require.Equal(t, w.Body.String(), `{"ans":"123"}`)
	})
}

func TestContext_QueryParam(t *testing.T) {
	r := httptest.NewRequest("GET", "http://example.com/foo/123?id=456&other=hello&boo=true", nil)
	w := httptest.NewRecorder()

	c := NewContext[any](w, r, readOptions{})

	t.Run("string", func(t *testing.T) {
		param := c.QueryParam("other")
		require.NotEmpty(t, param)
		require.Equal(t, param, "hello")

		param = c.QueryParam("notfound")
		require.Empty(t, param)
	})

	t.Run("int", func(t *testing.T) {
		param := c.QueryParam("id")
		require.NotEmpty(t, param)
		require.Equal(t, param, "456")

		paramInt := c.QueryParamInt("id", 0)
		require.Equal(t, paramInt, 456)

		paramInt = c.QueryParamInt("notfound", 42)
		require.Equal(t, paramInt, 42)

		paramInt = c.QueryParamInt("other", 42)
		require.Equal(t, paramInt, 42)

		paramInt, err := c.QueryParamIntErr("id")
		require.NoError(t, err)
		require.Equal(t, paramInt, 456)

		paramInt, err = c.QueryParamIntErr("notfound")
		require.Error(t, err)
		require.Equal(t, "param notfound not found", err.Error())
		require.Equal(t, 0, paramInt)

		paramInt, err = c.QueryParamIntErr("other")
		require.Error(t, err)
		require.Contains(t, err.Error(), "param other=hello is not of type int")
		require.Equal(t, paramInt, 0)
	})

	t.Run("bool", func(t *testing.T) {
		param := c.QueryParam("boo")
		require.NotEmpty(t, param)
		require.Equal(t, param, "true")

		paramBool := c.QueryParamBool("boo", false)
		require.Equal(t, paramBool, true)

		paramBool = c.QueryParamBool("notfound", true)
		require.Equal(t, paramBool, true)

		paramBool = c.QueryParamBool("other", true)
		require.Equal(t, paramBool, true)

		paramBool, err := c.QueryParamBoolErr("boo")
		require.NoError(t, err)
		require.Equal(t, paramBool, true)

		paramBool, err = c.QueryParamBoolErr("notfound")
		require.Error(t, err)
		require.Equal(t, paramBool, false)

		paramBool, err = c.QueryParamBoolErr("other")
		require.Error(t, err)
		require.Equal(t, paramBool, false)
	})
}

func TestContext_QueryParams(t *testing.T) {
	r := httptest.NewRequest("GET", "http://example.com/foo/123?id=456&other=hello", nil)
	w := httptest.NewRecorder()

	c := NewContext[any](w, r, readOptions{})

	params := c.QueryParams()
	require.NotEmpty(t, params)
	require.Equal(t, params["id"], "456")
	require.Equal(t, params["other"], "hello")
	require.Empty(t, params["notfound"])
}

type testStruct struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

type testStructInTransformer struct {
	Name string `json:"name" validate:"required,min=3,max=20"`
	Age  int    `json:"age" validate:"min=18"`
}

func (b *testStructInTransformer) InTransform(context.Context) error {
	b.Name = "transformed " + b.Name
	b.Age *= 2
	return nil
}

type testStructInTransformerWithError struct {
	Name string `json:"name" validate:"required,min=3,max=10"`
	Age  int    `json:"age" validate:"min=18"`
}

func (b *testStructInTransformerWithError) InTransform(context.Context) error {
	return errors.New("error")
}

func TestContext_Body(t *testing.T) {
	t.Run("can read JSON body", func(t *testing.T) {
		// Create new Reader
		a := strings.NewReader(`{"name":"John","age":30}`)

		// Test an http request
		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "http://example.com/foo", a)

		c := NewContext[testStruct](w, r, readOptions{})

		body, err := c.Body()
		require.NoError(t, err)
		require.Equal(t, body.Name, "John")
		require.Equal(t, body.Age, 30)
	})

	t.Run("can read JSON body twice", func(t *testing.T) {
		a := strings.NewReader(`{"name":"John","age":30}`)

		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "http://example.com/foo", a)

		c := NewContext[testStruct](w, r, readOptions{})

		body, err := c.Body()
		require.NoError(t, err)
		require.Equal(t, body.Name, "John")
		require.Equal(t, body.Age, 30)

		body, err = c.Body()
		require.NoError(t, err)
		require.Equal(t, body.Name, "John")
		require.Equal(t, body.Age, 30)
	})

	t.Run("can read and validate valid JSON body", func(t *testing.T) {
		type testStruct struct {
			Name string `json:"name" validate:"required,min=3,max=10"`
			Age  int    `json:"age" validate:"min=18"`
		}

		reqBody := strings.NewReader(`{"name":"John","age":30}`)
		c := NewContext[testStruct](
			httptest.NewRecorder(),
			httptest.NewRequest("GET", "http://example.com/foo", reqBody),
			readOptions{})

		body, err := c.Body()
		require.NoError(t, err)
		require.Equal(t, body.Name, "John")
		require.Equal(t, body.Age, 30)
	})

	t.Run("can read and validate invalid JSON body", func(t *testing.T) {
		type testStruct struct {
			Name string `json:"name" validate:"required,min=3,max=10"`
			Age  int    `json:"age" validate:"min=18"`
		}

		reqBody := strings.NewReader(`{"name":"VeryLongName","age":12}`)
		c := NewContext[testStruct](
			httptest.NewRecorder(),
			httptest.NewRequest("GET", "http://example.com/foo", reqBody),
			readOptions{})

		body, err := c.Body()
		require.Error(t, err)
		require.Equal(t, body.Name, "VeryLongName")
		require.Equal(t, body.Age, 12)
	})

	t.Run("can transform JSON body with custom method", func(t *testing.T) {
		reqBody := strings.NewReader(`{"name":"John","age":30}`)
		c := NewContext[testStructInTransformer](
			httptest.NewRecorder(),
			httptest.NewRequest("GET", "http://example.com/foo", reqBody),
			readOptions{})

		body, err := c.Body()
		require.NoError(t, err)
		require.Equal(t, body.Name, "transformed John")
		require.Equal(t, body.Age, 60)
	})

	t.Run("can transform JSON body with custom method returning error", func(t *testing.T) {
		reqBody := strings.NewReader(`{"name":"John","age":30}`)
		c := NewContext[testStructInTransformerWithError](
			httptest.NewRecorder(),
			httptest.NewRequest("GET", "http://example.com/foo", reqBody),
			readOptions{})

		body, err := c.Body()
		require.Error(t, err)
		require.Equal(t, body.Name, "John")
		require.Equal(t, body.Age, 30)
	})

	t.Run("unparsable because restricted to 1 byte", func(t *testing.T) {
		reqBody := strings.NewReader(`{"name":"John","age":30}`)
		c := NewContext[testStructInTransformerWithError](
			httptest.NewRecorder(),
			httptest.NewRequest("GET", "http://example.com/foo", reqBody),
			readOptions{
				MaxBodySize: 1,
			})

		body, err := c.Body()
		require.Error(t, err)
		require.Equal(t, "", body.Name)
		require.Equal(t, 0, body.Age)
	})

	t.Run("can read string body", func(t *testing.T) {
		// Create new Reader
		a := strings.NewReader("Hello World")

		// Test an http request
		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "http://example.com/foo", a)
		r.Header.Set("Content-Type", "text/plain")

		c := NewContext[string](w, r, readOptions{})

		_, err := c.Body()
		require.NoError(t, err)
	})
}

func FuzzContext_Body(f *testing.F) {
	f.Add("Hello Fuzz")

	f.Fuzz(func(t *testing.T, s string) {
		// Create new Reader
		a := strings.NewReader(s)

		// Test an http request
		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "http://example.com/foo", a)
		r.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		c := NewContext[testStruct](w, r, readOptions{})

		_, err := c.Body()
		require.NoError(t, err)
	})
}

func BenchmarkContext_Body(b *testing.B) {
	b.Run("valid JSON body", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			reqBody := strings.NewReader(`{"name":"John","age":30}`)
			c := NewContext[testStruct](
				httptest.NewRecorder(),
				httptest.NewRequest("GET", "http://example.com/foo", reqBody),
				readOptions{})
			_, err := c.Body()
			if err != nil {
				b.Fatal(err, "iteration", i)
			}
		}
	})

	// This test does make really sense because the body will not be accessed millions of times.
	// It however does show that caching the body works.
	// See [Body] for more information.
	b.Run("valid JSON body cache", func(b *testing.B) {
		reqBody := strings.NewReader(`{"name":"John","age":30}`)
		c := NewContext[testStruct](
			httptest.NewRecorder(),
			httptest.NewRequest("GET", "http://example.com/foo", reqBody),
			readOptions{})
		for i := 0; i < b.N; i++ {
			_, err := c.Body()
			if err != nil {
				b.Fatal(err, "iteration", i)
			}
		}
	})

	b.Run("invalid JSON body", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			reqBody := strings.NewReader(`{"name":"John","age":30}`)
			c := NewContext[testStruct](
				httptest.NewRecorder(),
				httptest.NewRequest("GET", "http://example.com/foo", reqBody),
				readOptions{})
			_, err := c.Body()
			if err != nil {
				b.Fatal(err)
			}
		}
	})

	b.Run("string body", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			reqBody := strings.NewReader(`{"name":"John","age":30}`)
			c := NewContext[testStruct](
				httptest.NewRecorder(),
				httptest.NewRequest("GET", "http://example.com/foo", reqBody),
				readOptions{})
			_, err := c.Body()
			if err != nil {
				b.Fatal(err)
			}
		}
	})
}

func TestContext_MustBody(t *testing.T) {
	t.Run("can read JSON body", func(t *testing.T) {
		// Create new Reader
		a := strings.NewReader(`{"name":"John","age":30}`)

		// Test an http request
		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "http://example.com/foo", a)

		c := NewContext[testStruct](w, r, readOptions{})

		body := c.MustBody()
		require.Equal(t, body.Name, "John")
		require.Equal(t, body.Age, 30)
	})

	t.Run("cannot read invalid JSON body", func(t *testing.T) {
		type testStruct struct {
			Name string `json:"name" validate:"required,min=3,max=10"`
			Age  int    `json:"age" validate:"min=18"`
		}

		reqBody := strings.NewReader(`{"name":"VeryLongName","age":12}`)
		c := NewContext[testStruct](
			httptest.NewRecorder(),
			httptest.NewRequest("GET", "http://example.com/foo", reqBody),
			readOptions{})

		require.Panics(t, func() {
			c.MustBody()
		})
	})
}

func TestMainLang(t *testing.T) {
	r := httptest.NewRequest("GET", "/", nil)
	r.Header.Set("Accept-Language", "fr-CH, fr;q=0.9, en;q=0.8, de;q=0.7, *;q=0.5")

	c := NewContext[any](httptest.NewRecorder(), r, readOptions{})
	require.Equal(t, c.MainLang(), "fr")
	require.Equal(t, c.MainLocale(), "fr-CH")
}

func TestContextNoBody_Body(t *testing.T) {
	body := `{"name":"John","age":30}`
	r := httptest.NewRequest("GET", "/", strings.NewReader(body))
	ctx := ContextNoBody{
		request:  r,
		response: httptest.NewRecorder(),
	}
	res, err := ctx.Body()
	require.NoError(t, err)
	require.Equal(t, any(map[string]any{
		"name": "John",
		"age":  30.0, // JSON numbers are float64
	}), res)
}

func TestContextNoBody_MustBody(t *testing.T) {
	t.Run("can read JSON body", func(t *testing.T) {
		body := `{"name":"John","age":30}`
		r := httptest.NewRequest("GET", "/", strings.NewReader(body))
		ctx := ContextNoBody{
			request:  r,
			response: httptest.NewRecorder(),
		}
		res := ctx.MustBody()
		require.Equal(t, any(map[string]any{
			"name": "John",
			"age":  30.0, // JSON numbers are float64
		}), res)
	})

	t.Run("cannot read invalid JSON body", func(t *testing.T) {
		body := `{"name":"John","age":30`
		r := httptest.NewRequest("GET", "/", strings.NewReader(body))
		ctx := ContextNoBody{
			request:  r,
			response: httptest.NewRecorder(),
		}
		require.Panics(t, func() {
			ctx.MustBody()
		})
	})
}

func TestContextNoBody_Redirect(t *testing.T) {
	s := NewServer()

	Get(s, "/", func(c ContextNoBody) (any, error) {
		return c.Redirect(301, "/foo")
	})

	Get(s, "/foo", func(c ContextNoBody) (ans, error) {
		return ans{Ans: "foo"}, nil
	})

	t.Run("can redirect", func(t *testing.T) {
		r := httptest.NewRequest("GET", "/", nil)
		w := httptest.NewRecorder()

		s.Mux.ServeHTTP(w, r)

		require.Equal(t, 301, w.Code)
		require.Equal(t, "/foo", w.Header().Get("Location"))
		require.Equal(t, "<a href=\"/foo\">Moved Permanently</a>.\n\n", w.Body.String())
	})
}
