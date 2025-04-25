package fuego

import (
	"bytes"
	"context"
	"encoding/xml"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/go-fuego/fuego/internal"
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
		Get(s, "/foo/{id}", func(c ContextNoBody) (any, error) {
			return c.PathParamIntErr("id")
		})

		r := httptest.NewRequest("GET", "/foo/abc", nil)
		w := httptest.NewRecorder()

		s.Mux.ServeHTTP(w, r)

		assert.Equal(t, http.StatusUnprocessableEntity, w.Code)
		assert.JSONEq(t, `{"title":"Unprocessable Entity","status":422,"detail":"path param id=abc is not of type int"}`, w.Body.String())
	})

	t.Run("path param not found", func(t *testing.T) {
		s := NewServer()
		Get(s, "/foo/", func(c ContextNoBody) (any, error) {
			return c.PathParamIntErr("id")
		})

		r := httptest.NewRequest("GET", "/foo/", nil)
		w := httptest.NewRecorder()

		s.Mux.ServeHTTP(w, r)

		assert.Equal(t, http.StatusUnprocessableEntity, w.Code)
		assert.JSONEq(t, `{"title":"Unprocessable Entity","status":422,"detail":"path param id not found"}`, w.Body.String())
	})
}

func TestContext_QueryParam(t *testing.T) {
	r := httptest.NewRequest("GET", "http://example.com/foo/123?id=456&other=hello&boo=true&name=jhon&name=doe", nil)
	w := httptest.NewRecorder()

	c := NewNetHTTPContext[any, any](BaseRoute{}, w, r, readOptions{})

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
		assert.Zero(t, paramInt)
		var notFoundErr internal.QueryParamNotFoundError
		require.ErrorAs(t, err, &notFoundErr)
		assert.Equal(t, "param notfound not found", notFoundErr.Error())
		assert.Equal(t, http.StatusUnprocessableEntity, notFoundErr.StatusCode())
		assert.Equal(t, "param notfound not found", notFoundErr.DetailMsg())

		paramInt, err = c.QueryParamIntErr("other")
		require.Error(t, err)
		require.Zero(t, paramInt)
		var invalidErr internal.QueryParamInvalidTypeError
		require.ErrorAs(t, err, &invalidErr)
		assert.Equal(t, `query param other=hello is not of type int: strconv.Atoi: parsing "hello": invalid syntax`, invalidErr.Error())
		assert.Equal(t, http.StatusUnprocessableEntity, invalidErr.StatusCode())
		assert.Equal(t, "query param other=hello is not of type int", invalidErr.DetailMsg())
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
		assert.False(t, paramBool)
		notFoundErr := &internal.QueryParamNotFoundError{}
		require.ErrorAs(t, err, notFoundErr)
		assert.Equal(t, "param notfound not found", notFoundErr.Error())
		assert.Equal(t, http.StatusUnprocessableEntity, notFoundErr.StatusCode())
		assert.Equal(t, "param notfound not found", notFoundErr.DetailMsg())

		paramBool, err = c.QueryParamBoolErr("other")
		require.Error(t, err)
		assert.False(t, paramBool)
		invalidErr := &internal.QueryParamInvalidTypeError{}
		require.ErrorAs(t, err, invalidErr)
		assert.Equal(t, `query param other=hello is not of type bool: strconv.ParseBool: parsing "hello": invalid syntax`, invalidErr.Error())
		assert.Equal(t, http.StatusUnprocessableEntity, invalidErr.StatusCode())
		assert.Equal(t, "query param other=hello is not of type bool", invalidErr.DetailMsg())
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

	c := NewNetHTTPContext[any, any](BaseRoute{}, w, r, readOptions{})

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

		c := NewNetHTTPContext[testStruct, any](BaseRoute{}, w, r, readOptions{})

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

		c := NewNetHTTPContext[testStruct, any](BaseRoute{}, w, r, readOptions{})

		body, err := c.Body()
		require.NoError(t, err)
		require.Equal(t, "John", body.Name)
		require.Equal(t, 30, body.Age)
	})

	t.Run("can read JSON body twice", func(t *testing.T) {
		a := strings.NewReader(`{"name":"John","age":30}`)

		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "http://example.com/foo", a)

		c := NewNetHTTPContext[testStruct, any](BaseRoute{}, w, r, readOptions{})

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
		c := NewNetHTTPContext[testStruct, any](
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
		c := NewNetHTTPContext[testStruct, any](
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
		c := NewNetHTTPContext[testStructInTransformer, any](
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
		c := NewNetHTTPContext[testStructInTransformerWithError, any](
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

		c := NewNetHTTPContext[[]byte, any](BaseRoute{}, w, r, readOptions{})
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

		c := NewNetHTTPContext[*struct{}, any](BaseRoute{}, w, r, readOptions{})
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

		c := NewNetHTTPContext[testStruct, any](BaseRoute{}, w, r, readOptions{})

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

		c := NewNetHTTPContext[testStruct, any](BaseRoute{}, w, r, readOptions{})

		body, err := c.Body()
		require.NoError(t, err)
		require.Equal(t, "John", body.Name)
		require.Equal(t, 30, body.Age)
	})

	t.Run("unparsable because restricted to 1 byte", func(t *testing.T) {
		reqBody := strings.NewReader(`{"name":"John","age":30}`)
		c := NewNetHTTPContext[testStructInTransformerWithError, any](
			BaseRoute{}, httptest.NewRecorder(),
			httptest.NewRequest("GET", "http://example.com/foo", reqBody),
			readOptions{
				MaxBodySize: 1,
			})

		body, err := c.Body()
		require.Error(t, err)
		require.Empty(t, body.Name)
		require.Zero(t, body.Age)
	})

	t.Run("can read string body", func(t *testing.T) {
		// Create new Reader
		a := strings.NewReader("Hello World")

		// Test an http request
		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "http://example.com/foo", a)
		r.Header.Set("Content-Type", "text/plain")

		c := NewNetHTTPContext[string, any](BaseRoute{}, w, r, readOptions{})

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

		c := NewNetHTTPContext[testStruct, any](BaseRoute{}, w, r, readOptions{})

		_, err := c.Body()
		require.NoError(t, err)
	})
}

func BenchmarkContext_Body(b *testing.B) {
	b.Run("valid JSON body", func(b *testing.B) {
		for i := range b.N {
			reqBody := strings.NewReader(`{"name":"John","age":30}`)
			c := NewNetHTTPContext[testStruct, any](
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
		c := NewNetHTTPContext[testStruct, any](
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
			c := NewNetHTTPContext[testStruct, any](
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
			c := NewNetHTTPContext[testStruct, any](
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

		c := NewNetHTTPContext[testStruct, any](BaseRoute{}, w, r, readOptions{})

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
		c := NewNetHTTPContext[testStruct, any](
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

	c := NewNetHTTPContext[any, any](BaseRoute{}, httptest.NewRecorder(), r, readOptions{})
	assert.Equal(t, "fr", c.MainLang())
	require.Equal(t, "fr-CH", c.MainLocale())
}

func TestContextNoBody_Body(t *testing.T) {
	body := `{"name":"John","age":30}`
	r := httptest.NewRequest("GET", "/", strings.NewReader(body))
	ctx := netHttpContext[any, any]{
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
		ctx := netHttpContext[any, any]{
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
		ctx := netHttpContext[any, any]{
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

func TestNetHttpContext_Params(t *testing.T) {
	t.Run("can write and read params", func(t *testing.T) {
		type MyParams struct {
			ID          int     `query:"id"`
			Temperature float64 `query:"temperature"`
			Other       string  `query:"other" description:"my description"`
			ContentType string  `header:"Content-Type"`
		}
		r := httptest.NewRequest("GET", "http://example.com/foo/123?id=456&other=hello&temperature=20.30", nil)
		r.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		c := NewNetHTTPContext[any, MyParams](BaseRoute{}, w, r, readOptions{})

		params, err := c.Params()
		require.NoError(t, err)
		require.NotEmpty(t, params)
		assert.Equal(t, 456, params.ID)
		assert.Equal(t, "hello", params.Other)
		assert.Equal(t, "application/json", params.ContentType)
		assert.InEpsilon(t, 20.30, params.Temperature, 0.01)
	})

	t.Run("does not support other receivers than struct", func(t *testing.T) {
		t.Run("pointer to struct", func(t *testing.T) {
			type MyParams struct{}
			r := httptest.NewRequest("GET", "http://example.com/foo/123?id=456&other=hello&temperature=20.30", nil)
			r.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			c := NewNetHTTPContext[any, *MyParams](BaseRoute{}, w, r, readOptions{})

			_, err := c.Params()

			require.ErrorContains(t, err, "params must be a struct, got *fuego.MyParams")
		})

		t.Run("interface", func(t *testing.T) {
			r := httptest.NewRequest("GET", "http://example.com/foo/123?id=456&other=hello&temperature=20.30", nil)
			r.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			c := NewNetHTTPContext[any, any](BaseRoute{}, w, r, readOptions{})

			_, err := c.Params()

			require.ErrorContains(t, err, "params must be a struct, got <nil>")
		})
	})

	t.Run("support for more integer types", func(t *testing.T) {
		type MyParams struct {
			ID          int8    `query:"id"`
			Temperature float32 `query:"temperature"`
			Other       uint32  `query:"other" description:"my description"`
			MyHeader    uint64  `header:"MyHeader"`
		}

		r := httptest.NewRequest("GET", "http://example.com/foo/123?id=12&other=23782&temperature=20.30", nil)
		r.Header.Set("MyHeader", "8923")
		w := httptest.NewRecorder()
		c := NewNetHTTPContext[any, MyParams](BaseRoute{}, w, r, readOptions{})
		params, err := c.Params()
		require.NoError(t, err)
		require.NotEmpty(t, params)
		assert.Equal(t, int8(12), params.ID)
		assert.Equal(t, uint32(23782), params.Other)
		assert.Equal(t, uint64(8923), params.MyHeader)
		assert.InEpsilon(t, float32(20.30), params.Temperature, 0.01)
	})

	t.Run("support for array of strings", func(t *testing.T) {
		type MyParams struct {
			Tags []string `query:"tags"`
		}

		r := httptest.NewRequest("GET", "http://example.com/foo?tags=golang&tags=web&tags=api", nil)
		w := httptest.NewRecorder()
		c := NewNetHTTPContext[any, MyParams](BaseRoute{}, w, r, readOptions{})
		params, err := c.Params()
		require.NoError(t, err)
		require.NotEmpty(t, params)
		assert.Equal(t, []string{"golang", "web", "api"}, params.Tags)
	})

	t.Run("support for array of integers", func(t *testing.T) {
		type MyParams struct {
			IDs []int `query:"ids"`
		}

		r := httptest.NewRequest("GET", "http://example.com/foo?ids=1&ids=2&ids=3", nil)
		w := httptest.NewRecorder()
		c := NewNetHTTPContext[any, MyParams](BaseRoute{}, w, r, readOptions{})
		params, err := c.Params()
		require.NoError(t, err)
		require.NotEmpty(t, params)
		assert.Equal(t, []int{1, 2, 3}, params.IDs)
	})

	t.Run("support for array of various integer types", func(t *testing.T) {
		type MyParams struct {
			Int8s   []int8   `query:"int8s"`
			Int16s  []int16  `query:"int16s"`
			Int32s  []int32  `query:"int32s"`
			Int64s  []int64  `query:"int64s"`
			Uints   []uint   `query:"uints"`
			Uint8s  []uint8  `query:"uint8s"`
			Uint16s []uint16 `query:"uint16s"`
			Uint32s []uint32 `query:"uint32s"`
			Uint64s []uint64 `query:"uint64s"`
		}

		url := "http://example.com/foo?" +
			"int8s=1&int8s=2&" +
			"int16s=300&int16s=400&" +
			"int32s=70000&int32s=80000&" +
			"int64s=9000000000&int64s=9100000000&" +
			"uints=1&uints=2&" +
			"uint8s=200&uint8s=201&" +
			"uint16s=60000&uint16s=60001&" +
			"uint32s=3000000000&uint32s=3100000000&" +
			"uint64s=9000000000&uint64s=9100000000"

		r := httptest.NewRequest("GET", url, nil)
		w := httptest.NewRecorder()
		c := NewNetHTTPContext[any, MyParams](BaseRoute{}, w, r, readOptions{})
		params, err := c.Params()
		require.NoError(t, err)
		require.NotEmpty(t, params)
		assert.Equal(t, []int8{1, 2}, params.Int8s)
		assert.Equal(t, []int16{300, 400}, params.Int16s)
		assert.Equal(t, []int32{70000, 80000}, params.Int32s)
		assert.Equal(t, []int64{9000000000, 9100000000}, params.Int64s)
		assert.Equal(t, []uint{1, 2}, params.Uints)
		assert.Equal(t, []uint8{200, 201}, params.Uint8s)
		assert.Equal(t, []uint16{60000, 60001}, params.Uint16s)
		assert.Equal(t, []uint32{3000000000, 3100000000}, params.Uint32s)
		assert.Equal(t, []uint64{9000000000, 9100000000}, params.Uint64s)
	})

	t.Run("support for array of booleans", func(t *testing.T) {
		type MyParams struct {
			Flags []bool `query:"flags"`
		}

		r := httptest.NewRequest("GET", "http://example.com/foo?flags=true&flags=false&flags=true", nil)
		w := httptest.NewRecorder()
		c := NewNetHTTPContext[any, MyParams](BaseRoute{}, w, r, readOptions{})
		params, err := c.Params()
		require.NoError(t, err)
		require.NotEmpty(t, params)
		assert.Equal(t, []bool{true, false, true}, params.Flags)
	})

	t.Run("support for array of floats", func(t *testing.T) {
		type MyParams struct {
			Float32s []float32 `query:"float32s"`
			Float64s []float64 `query:"float64s"`
		}

		r := httptest.NewRequest("GET", "http://example.com/foo?float32s=1.1&float32s=-2.2&float64s=-3.3&float64s=4.4", nil)
		w := httptest.NewRecorder()
		c := NewNetHTTPContext[any, MyParams](BaseRoute{}, w, r, readOptions{})
		params, err := c.Params()
		require.NoError(t, err)
		require.NotEmpty(t, params)
		assert.InEpsilonSlice(t, []float32{1.1, -2.2}, params.Float32s, 0.01)
		assert.InEpsilonSlice(t, []float64{-3.3, 4.4}, params.Float64s, 0.01)
	})

	t.Run("error handling for invalid array values", func(t *testing.T) {
		type MyParams struct {
			IDs []int `query:"ids"`
		}

		r := httptest.NewRequest("GET", "http://example.com/foo?ids=1&ids=invalid&ids=3", nil)
		w := httptest.NewRecorder()
		c := NewNetHTTPContext[any, MyParams](BaseRoute{}, w, r, readOptions{})
		_, err := c.Params()
		require.Error(t, err)
		assert.Contains(t, err.Error(), "cannot convert invalid to int")
	})

	t.Run("empty array when no query parameters", func(t *testing.T) {
		type MyParams struct {
			Tags []string `query:"tags"`
		}

		r := httptest.NewRequest("GET", "http://example.com/foo", nil)
		w := httptest.NewRecorder()
		c := NewNetHTTPContext[any, MyParams](BaseRoute{}, w, r, readOptions{})
		params, err := c.Params()
		require.NoError(t, err)
		assert.Empty(t, params.Tags)
	})

	t.Run("mixed single and array parameters", func(t *testing.T) {
		type MyParams struct {
			ID    int      `query:"id"`
			Tags  []string `query:"tags"`
			Limit int      `query:"limit"`
		}

		r := httptest.NewRequest("GET", "http://example.com/foo?id=123&tags=golang&tags=web&limit=50", nil)
		w := httptest.NewRecorder()
		c := NewNetHTTPContext[any, MyParams](BaseRoute{}, w, r, readOptions{})
		params, err := c.Params()
		require.NoError(t, err)
		require.NotEmpty(t, params)
		assert.Equal(t, 123, params.ID)
		assert.Equal(t, []string{"golang", "web"}, params.Tags)
		assert.Equal(t, 50, params.Limit)
	})

	t.Run("error handling for integer overflow", func(t *testing.T) {
		t.Run("int8 overflow", func(t *testing.T) {
			r := httptest.NewRequest("GET", "http://example.com/foo?value=128", nil)
			w := httptest.NewRecorder()

			ctx := NewNetHTTPContext[any, struct {
				Value int8 `query:"value"`
			}](BaseRoute{}, w, r, readOptions{})

			_, err := ctx.Params()
			require.Error(t, err)
			assert.Contains(t, err.Error(), "cannot convert 128 to int")
		})

		t.Run("int8 underflow", func(t *testing.T) {
			r := httptest.NewRequest("GET", "http://example.com/foo?value=-129", nil)
			w := httptest.NewRecorder()

			ctx := NewNetHTTPContext[any, struct {
				Value int8 `query:"value"`
			}](BaseRoute{}, w, r, readOptions{})

			_, err := ctx.Params()
			require.Error(t, err)
			assert.Contains(t, err.Error(), "cannot convert -129 to int")
		})

		t.Run("int16 overflow", func(t *testing.T) {
			r := httptest.NewRequest("GET", "http://example.com/foo?value=32768", nil)
			w := httptest.NewRecorder()

			ctx := NewNetHTTPContext[any, struct {
				Value int16 `query:"value"`
			}](BaseRoute{}, w, r, readOptions{})

			_, err := ctx.Params()
			require.Error(t, err)
			assert.Contains(t, err.Error(), "cannot convert 32768 to int")
		})

		t.Run("uint8 overflow", func(t *testing.T) {
			r := httptest.NewRequest("GET", "http://example.com/foo?value=256", nil)
			w := httptest.NewRecorder()

			ctx := NewNetHTTPContext[any, struct {
				Value uint8 `query:"value"`
			}](BaseRoute{}, w, r, readOptions{})

			_, err := ctx.Params()
			require.Error(t, err)
			assert.Contains(t, err.Error(), "cannot convert 256 to uint")
		})

		t.Run("negative value for unsigned type", func(t *testing.T) {
			r := httptest.NewRequest("GET", "http://example.com/foo?value=-1", nil)
			w := httptest.NewRecorder()

			ctx := NewNetHTTPContext[any, struct {
				Value uint32 `query:"value"`
			}](BaseRoute{}, w, r, readOptions{})

			_, err := ctx.Params()
			require.Error(t, err)
			assert.Contains(t, err.Error(), "cannot convert -1 to uint")
		})

		t.Run("array with overflow value", func(t *testing.T) {
			r := httptest.NewRequest("GET", "http://example.com/foo?values=100&values=300", nil)
			w := httptest.NewRecorder()

			ctx := NewNetHTTPContext[any, struct {
				Values []uint8 `query:"values"`
			}](BaseRoute{}, w, r, readOptions{})

			_, err := ctx.Params()
			require.Error(t, err)
			assert.Contains(t, err.Error(), "cannot convert 300 to uint")
		})

		t.Run("valid values within range", func(t *testing.T) {
			r := httptest.NewRequest("GET", "http://example.com/foo?int8=127&uint8=255&int16=32767&uint16=65535", nil)
			w := httptest.NewRecorder()

			ctx := NewNetHTTPContext[any, struct {
				Int8   int8   `query:"int8"`
				Uint8  uint8  `query:"uint8"`
				Int16  int16  `query:"int16"`
				Uint16 uint16 `query:"uint16"`
			}](BaseRoute{}, w, r, readOptions{})

			_, err := ctx.Params()
			require.NoError(t, err)
		})
	})
}
