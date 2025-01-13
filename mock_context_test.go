package fuego

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

type TestBody struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func TestMockContext(t *testing.T) {
	// Create a new mock context
	ctx := NewMockContext[TestBody]()

	// Test body
	body := TestBody{
		Name: "John",
		Age:  30,
	}
	ctx.SetBody(body)
	gotBody, err := ctx.Body()
	assert.NoError(t, err)
	assert.Equal(t, body, gotBody)

	// Test URL values
	values := url.Values{
		"key": []string{"value"},
	}
	ctx.SetURLValues(values)
	assert.Equal(t, values, ctx.URLValues())

	// Test headers
	ctx.SetHeader("Content-Type", "application/json")
	assert.Equal(t, "application/json", ctx.Header().Get("Content-Type"))

	// Test path params
	ctx.SetPathParam("id", "123")
	assert.Equal(t, "123", ctx.PathParam("id"))
}

func TestMockContextAdvanced(t *testing.T) {
	// Test with custom context
	ctx := NewMockContext[TestBody]()
	customCtx := context.WithValue(context.Background(), "key", "value")
	ctx.SetContext(customCtx)
	assert.Equal(t, "value", ctx.Context().Value("key"))

	// Test with request/response
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/test", nil)
	ctx.SetResponse(w)
	ctx.SetRequest(r)
	assert.Equal(t, w, ctx.Response())
	assert.Equal(t, r, ctx.Request())

	// Test multiple headers
	ctx.SetHeader("X-Test-1", "value1")
	ctx.SetHeader("X-Test-2", "value2")
	assert.Equal(t, "value1", ctx.Header().Get("X-Test-1"))
	assert.Equal(t, "value2", ctx.Header().Get("X-Test-2"))

	// Test multiple path params
	ctx.SetPathParam("id", "123")
	ctx.SetPathParam("category", "books")
	assert.Equal(t, "123", ctx.PathParam("id"))
	assert.Equal(t, "books", ctx.PathParam("category"))
} 