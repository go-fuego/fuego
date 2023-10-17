package op

import (
	"errors"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

type testStructNormalizable struct {
	Name string `json:"name" validate:"required,min=3,max=10"`
	Age  int    `json:"age" validate:"min=18"`
}

func (b *testStructNormalizable) Normalize() error {
	b.Name = "normalized " + b.Name
	b.Age *= 2
	return nil
}

type testStructNormalizableWithError struct {
	Name string `json:"name" validate:"required,min=3,max=10"`
	Age  int    `json:"age" validate:"min=18"`
}

func (b *testStructNormalizableWithError) Normalize() error {
	return errors.New("error")
}

func TestContext_Body(t *testing.T) {
	t.Run("can read JSON body", func(t *testing.T) {
		type testStruct struct {
			Name string `json:"name"`
			Age  int    `json:"age"`
		}

		// Create new Reader
		a := strings.NewReader(`{"name":"John","age":30}`)

		// Test an http request
		r := httptest.NewRequest("GET", "http://example.com/foo", a)

		c := &Context[testStruct]{
			request: r,
		}

		body, err := c.Body()
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
		c := &Context[testStruct]{
			request: httptest.NewRequest("GET", "http://example.com/foo", reqBody),
		}
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
		c := &Context[testStruct]{
			request: httptest.NewRequest("GET", "http://example.com/foo", reqBody),
		}
		body, err := c.Body()
		require.Error(t, err)
		require.Equal(t, body.Name, "VeryLongName")
		require.Equal(t, body.Age, 12)
	})

	t.Run("can normalize JSON body with custom method", func(t *testing.T) {
		reqBody := strings.NewReader(`{"name":"John","age":30}`)
		c := &Context[testStructNormalizable]{
			request: httptest.NewRequest("GET", "http://example.com/foo", reqBody),
		}
		body, err := c.Body()
		require.NoError(t, err)
		require.Equal(t, body.Name, "normalized John")
		require.Equal(t, body.Age, 60)
	})

	t.Run("can normalize JSON body with custom method returning error", func(t *testing.T) {
		reqBody := strings.NewReader(`{"name":"John","age":30}`)
		c := &Context[testStructNormalizableWithError]{
			request: httptest.NewRequest("GET", "http://example.com/foo", reqBody),
		}
		body, err := c.Body()
		require.Error(t, err)
		require.Equal(t, body.Name, "John")
		require.Equal(t, body.Age, 30)
	})

	t.Run("can read string body", func(t *testing.T) {
		// Create new Reader
		a := strings.NewReader("Hello World")

		// Test an http request
		r := httptest.NewRequest("GET", "http://example.com/foo", a)

		c := &Context[string]{
			request: r,
		}

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
		r := httptest.NewRequest("GET", "http://example.com/foo", a)

		c := &Context[string]{
			request: r,
		}

		_, err := c.Body()
		require.NoError(t, err)
	})
}
