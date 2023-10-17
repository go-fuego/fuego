package op

import (
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

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

		c.Body()
	})
}
