package fuego

import (
	"bytes"
	"context"
	"encoding/xml"
	"errors"
	"fmt"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestContext_PathParam(t *testing.T) {
	t.Run("can read one path param", func(t *testing.T) {
		s := NewServer()
		Get(s, "/foo/{id}", func(c ContextNoBody) (ans, error) {
			return ans{Ans: c.PathParam("id")}, nil
		})

		r := httptest.NewRequest("GET", "/foo/123", nil)
		w := httptest.NewRecorder()

		s.Mux.ServeHTTP(w, r)

		require.Equal(t, crlf(`{"ans":"123"}`), w.Body.String())
	})

	t.Run("can read one path param to int", func(t *testing.T) {
		s := NewServer()
		Get(s, "/foo/{id}", func(c ContextNoBody) (ans, error) {
			return ans{Ans: fmt.Sprintf("%d", c.PathParamInt("id"))}, nil
		})

		r := httptest.NewRequest("GET", "/foo/123", nil)
		w := httptest.NewRecorder()

		s.Mux.ServeHTTP(w, r)

		require.Equal(t, crlf(`{"ans":"123"}`), w.Body.String())
	})

	t.Run("reading non-int path param to int defaults to 0", func(t *testing.T) {
		s := NewServer()
		Get(s, "/foo/{id}", func(c ContextNoBody) (ans, error) {
			return ans{Ans: fmt.Sprintf("%d", c.PathParamInt("id"))}, nil
		})

		r := httptest.NewRequest("GET", "/foo/abc", nil)
		w := httptest.NewRecorder()

		s.Mux.ServeHTTP(w, r)

		require.Equal(t, crlf(`{"ans":"0"}`), w.Body.String())
	})

	t.Run("reading missing path param to int defaults to 0", func(t *testing.T) {
		s := NewServer()
		Get(s, "/foo/", func(c ContextNoBody) (ans, error) {
			return ans{Ans: fmt.Sprintf("%d", c.PathParamInt("id"))}, nil
		})

		r := httptest.NewRequest("GET", "/foo/", nil)
		w := httptest.NewRecorder()

		s.Mux.ServeHTTP(w, r)

		require.Equal(t, crlf(`{"ans":"0"}`), w.Body.String())
	})

	t.Run("reading non-int path param to int sends an error", func(t *testing.T) {
		s := NewServer()
		Get(s, "/foo/{id}", func(c ContextNoBody) (ans, error) {
			id, err := c.PathParamIntErr("id")
			if err != nil {
				return ans{}, err
			}
			return ans{Ans: fmt.Sprintf("%d", id)}, nil
		})

		r := httptest.NewRequest("GET", "/foo/abc", nil)
		w := httptest.NewRecorder()

		s.Mux.ServeHTTP(w, r)

		require.Equal(t, 422, w.Code)
	})

	t.Run("path param invalid", func(t *testing.T) {
		s := NewServer()
		Get(s, "/foo/", func(c ContextNoBody) (ans, error) {
			return ans{Ans: c.PathParam("id")}, nil
		})

		r := httptest.NewRequest("GET", "/foo/", nil)
		w := httptest.NewRecorder()

		s.Mux.ServeHTTP(w, r)

		require.Equal(t, crlf(`{"ans":""}`), w.Body.String())
	})
}

func TestContext_QueryParam(t *testing.T) {
	r := httptest.NewRequest("GET", "http://example.com/foo/123?id=456&other=hello&boo=true&name=jhon&name=doe", nil)
	w := httptest.NewRecorder()

	c := NewNetHTTPContext[any](BaseRoute{}, w, r, readOptions{})

	t.Run("string", func(t *testing.T) {
		param := c.QueryParam("other")
		require.NotEmpty(t, param)
		require.Equal(t, "hello", param)

		param = c.QueryParam("notfound")
		require.Empty(t, param)
	})

	t.Run("int", func(t *testing.T) {
		param := c.QueryParam("id")
		require.NotEmpty(t, param)
		require.Equal(t, "456", param)

		paramInt := c.QueryParamInt("id")
		require.Equal(t, 456, paramInt)

		paramInt = c.QueryParamInt("notfound")
		require.Zero(t, paramInt)

		paramInt = c.QueryParamInt("other")
		require.Zero(t, paramInt)

		paramInt, err := c.QueryParamIntErr("id")
		require.NoError(t, err)
		require.Equal(t, 456, paramInt)

		paramInt, err = c.QueryParamIntErr("notfound")
		require.Error(t, err)
		require.Equal(t, "param notfound not found", err.Error())
		require.Zero(t, paramInt)

		paramInt, err = c.QueryParamIntErr("other")
		require.Error(t, err)
		require.Contains(t, err.Error(), "param other=hello is not of type int")
		require.Zero(t, paramInt)
	})

	t.Run("bool", func(t *testing.T) {
		param := c.QueryParam("boo")
		require.NotEmpty(t, param)
		require.Equal(t, "true", param)

		paramBool := c.QueryParamBool("boo")
		require.True(t, paramBool)

		paramBool = c.QueryParamBool("notfound")
		require.False(t, paramBool)

		paramBool, err := c.QueryParamBoolErr("boo")
		require.NoError(t, err)
		require.True(t, paramBool)

		paramBool, err = c.QueryParamBoolErr("notfound")
		require.Error(t, err)
		require.False(t, paramBool)

		paramBool, err = c.QueryParamBoolErr("other")
		require.Error(t, err)
		require.False(t, paramBool)
	})

	t.Run("slice", func(t *testing.T) {
		name := c.QueryParamArr("name")
		require.NotEmpty(t, name)
		require.Equal(t, []string{"jhon", "doe"}, name)

		notFound := c.QueryParamArr("notfound")
		require.Empty(t, notFound)
	})
}

func TestContext_QueryParams(t *testing.T) {
	r := httptest.NewRequest("GET", "http://example.com/foo/123?id=456&other=hello", nil)
	w := httptest.NewRecorder()

	c := NewNetHTTPContext[any](BaseRoute{}, w, r, readOptions{})

	params := c.QueryParams()
	require.NotEmpty(t, params)
	require.Equal(t, []string{"456"}, params["id"])
	require.Equal(t, []string{"hello"}, params["other"])
	require.Empty(t, params["notfound"])
}

type testStruct struct {
	XMLName xml.Name `xml:"TestStruct"`
	Name    string   `json:"name" xml:"Name" yaml:"name"`
	Age     int      `json:"age" xml:"Age" yaml:"age"`
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

		c := NewNetHTTPContext[testStruct](BaseRoute{}, w, r, readOptions{})

		body, err := c.Body()
		require.NoError(t, err)
		require.Equal(t, "John", body.Name)
		require.Equal(t, 30, body.Age)
	})

	t.Run("can read JSON body with Content-Type application/json", func(t *testing.T) {
		// Create new Reader
		a := strings.NewReader(`{"name":"John","age":30}`)

		// Test an http request
		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "http://example.com/foo", a)
		r.Header.Add("Content-Type", "application/json")

		c := NewNetHTTPContext[testStruct](BaseRoute{}, w, r, readOptions{})

		body, err := c.Body()
		require.NoError(t, err)
		require.Equal(t, "John", body.Name)
		require.Equal(t, 30, body.Age)
	})

	t.Run("can read JSON body twice", func(t *testing.T) {
		a := strings.NewReader(`{"name":"John","age":30}`)

		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "http://example.com/foo", a)

		c := NewNetHTTPContext[testStruct](BaseRoute{}, w, r, readOptions{})

		body, err := c.Body()
		require.NoError(t, err)
		require.Equal(t, "John", body.Name)
		require.Equal(t, 30, body.Age)

		body, err = c.Body()
		require.NoError(t, err)
		require.Equal(t, "John", body.Name)
		require.Equal(t, 30, body.Age)
	})

	t.Run("can read and validate valid JSON body", func(t *testing.T) {
		type testStruct struct {
			Name string `json:"name" validate:"required,min=3,max=10"`
			Age  int    `json:"age" validate:"min=18"`
		}

		reqBody := strings.NewReader(`{"name":"John","age":30}`)
		c := NewNetHTTPContext[testStruct](
			BaseRoute{},
			httptest.NewRecorder(),
			httptest.NewRequest("GET", "http://example.com/foo", reqBody),
			readOptions{})

		body, err := c.Body()
		require.NoError(t, err)
		require.Equal(t, "John", body.Name)
		require.Equal(t, 30, body.Age)
	})

	t.Run("can read and validate invalid JSON body", func(t *testing.T) {
		type testStruct struct {
			Name string `json:"name" validate:"required,min=3,max=10"`
			Age  int    `json:"age" validate:"min=18"`
		}

		reqBody := strings.NewReader(`{"name":"VeryLongName","age":12}`)
		c := NewNetHTTPContext[testStruct](
			BaseRoute{},
			httptest.NewRecorder(),
			httptest.NewRequest("GET", "http://example.com/foo", reqBody),
			readOptions{})

		body, err := c.Body()
		require.Error(t, err)
		require.Equal(t, "VeryLongName", body.Name)
		require.Equal(t, 12, body.Age)
	})

	t.Run("can transform JSON body with custom method", func(t *testing.T) {
		reqBody := strings.NewReader(`{"name":"John","age":30}`)
		c := NewNetHTTPContext[testStructInTransformer](
			BaseRoute{},
			httptest.NewRecorder(),
			httptest.NewRequest("GET", "http://example.com/foo", reqBody),
			readOptions{})

		body, err := c.Body()
		require.NoError(t, err)
		require.Equal(t, "transformed John", body.Name)
		require.Equal(t, 60, body.Age)
	})

	t.Run("can transform JSON body with custom method returning error", func(t *testing.T) {
		reqBody := strings.NewReader(`{"name":"John","age":30}`)
		c := NewNetHTTPContext[testStructInTransformerWithError](
			BaseRoute{}, httptest.NewRecorder(),
			httptest.NewRequest("GET", "http://example.com/foo", reqBody),
			readOptions{})

		body, err := c.Body()
		require.Error(t, err)
		require.Equal(t, "John", body.Name)
		require.Equal(t, 30, body.Age)
	})

	t.Run("can read bytes", func(t *testing.T) {
		// Create new Reader with pure bytes from an image
		a := bytes.NewReader([]byte(`image`))

		// Test an http request
		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "http://example.com/foo", a)
		r.Header.Add("Content-Type", "application/octet-stream")

		c := NewNetHTTPContext[[]byte](BaseRoute{}, w, r, readOptions{})
		body, err := c.Body()
		require.NoError(t, err)
		require.Equal(t, []byte(`image`), body)
	})

	t.Run("cannot read bytes if expected type is different than bytes", func(t *testing.T) {
		// Create new Reader with pure bytes from an image
		a := bytes.NewReader([]byte(`image`))

		// Test an http request
		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "http://example.com/foo", a)
		r.Header.Add("Content-Type", "application/octet-stream")

		c := NewNetHTTPContext[*struct{}](BaseRoute{}, w, r, readOptions{})
		body, err := c.Body()
		require.Error(t, err)
		require.ErrorContains(t, err, "use []byte as the body type")
		require.Equal(t, (*struct{})(nil), body)
	})

	t.Run("can read XML body", func(t *testing.T) {
		a := bytes.NewReader([]byte(`
<TestStruct>
	<Name>John</Name>
	<Age>30</Age>
</TestStruct>
`))
		// Test an http request
		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "http://example.com/foo", a)
		r.Header.Add("Content-Type", "application/xml")

		c := NewNetHTTPContext[testStruct](BaseRoute{}, w, r, readOptions{})

		body, err := c.Body()
		require.NoError(t, err)
		require.Equal(t, "John", body.Name)
		require.Equal(t, 30, body.Age)
	})

	t.Run("can read YAML body", func(t *testing.T) {
		a := bytes.NewReader([]byte(`
name: John
age: 30
`))
		// Test an http request
		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "http://example.com/foo", a)
		r.Header.Add("Content-Type", "application/x-yaml")

		c := NewNetHTTPContext[testStruct](BaseRoute{}, w, r, readOptions{})

		body, err := c.Body()
		require.NoError(t, err)
		require.Equal(t, "John", body.Name)
		require.Equal(t, 30, body.Age)
	})

	t.Run("unparsable because restricted to 1 byte", func(t *testing.T) {
		reqBody := strings.NewReader(`{"name":"John","age":30}`)
		c := NewNetHTTPContext[testStructInTransformerWithError](
			BaseRoute{}, httptest.NewRecorder(),
			httptest.NewRequest("GET", "http://example.com/foo", reqBody),
			readOptions{
				MaxBodySize: 1,
			})

		body, err := c.Body()
		require.Error(t, err)
		require.Equal(t, "", body.Name)
		require.Zero(t, body.Age)
	})

	t.Run("can read string body", func(t *testing.T) {
		// Create new Reader
		a := strings.NewReader("Hello World")

		// Test an http request
		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "http://example.com/foo", a)
		r.Header.Set("Content-Type", "text/plain")

		c := NewNetHTTPContext[string](BaseRoute{}, w, r, readOptions{})

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

		c := NewNetHTTPContext[testStruct](BaseRoute{}, w, r, readOptions{})

		_, err := c.Body()
		require.NoError(t, err)
	})
}

func BenchmarkContext_Body(b *testing.B) {
	b.Run("valid JSON body", func(b *testing.B) {
		for i := range b.N {
			reqBody := strings.NewReader(`{"name":"John","age":30}`)
			c := NewNetHTTPContext[testStruct](
				BaseRoute{}, httptest.NewRecorder(),
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
		c := NewNetHTTPContext[testStruct](
			BaseRoute{}, httptest.NewRecorder(),
			httptest.NewRequest("GET", "http://example.com/foo", reqBody),
			readOptions{})
		for i := range b.N {
			_, err := c.Body()
			if err != nil {
				b.Fatal(err, "iteration", i)
			}
		}
	})

	b.Run("invalid JSON body", func(b *testing.B) {
		for range b.N {
			reqBody := strings.NewReader(`{"name":"John","age":30}`)
			c := NewNetHTTPContext[testStruct](
				BaseRoute{}, httptest.NewRecorder(),
				httptest.NewRequest("GET", "http://example.com/foo", reqBody),
				readOptions{})
			_, err := c.Body()
			if err != nil {
				b.Fatal(err)
			}
		}
	})

	b.Run("string body", func(b *testing.B) {
		for range b.N {
			reqBody := strings.NewReader(`{"name":"John","age":30}`)
			c := NewNetHTTPContext[testStruct](
				BaseRoute{}, httptest.NewRecorder(),
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

		c := NewNetHTTPContext[testStruct](BaseRoute{}, w, r, readOptions{})

		body := c.MustBody()
		require.Equal(t, "John", body.Name)
		require.Equal(t, 30, body.Age)
	})

	t.Run("cannot read invalid JSON body", func(t *testing.T) {
		type testStruct struct {
			Name string `json:"name" validate:"required,min=3,max=10"`
			Age  int    `json:"age" validate:"min=18"`
		}

		reqBody := strings.NewReader(`{"name":"VeryLongName","age":12}`)
		c := NewNetHTTPContext[testStruct](
			BaseRoute{}, httptest.NewRecorder(),
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

	c := NewNetHTTPContext[any](BaseRoute{}, httptest.NewRecorder(), r, readOptions{})
	require.Equal(t, "fr", c.MainLang())
	require.Equal(t, "fr-CH", c.MainLocale())
}

func TestContextNoBody_Body(t *testing.T) {
	body := `{"name":"John","age":30}`
	r := httptest.NewRequest("GET", "/", strings.NewReader(body))
	ctx := netHttpContext[any]{
		Req: r,
		Res: httptest.NewRecorder(),
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
		ctx := netHttpContext[any]{
			Req: r,
			Res: httptest.NewRecorder(),
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
		ctx := netHttpContext[any]{
			Req: r,
			Res: httptest.NewRecorder(),
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
