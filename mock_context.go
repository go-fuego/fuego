package fuego

import (
	"context"
	"net/http"
	"net/url"
)

// MockContext provides a framework-agnostic implementation of ContextWithBody
// for testing purposes. It allows testing controllers without depending on
// specific web frameworks like Gin or Echo.
type MockContext[B any] struct {
	body       B
	urlValues  url.Values
	headers    http.Header
	pathParams map[string]string
	ctx        context.Context
	response   http.ResponseWriter
	request    *http.Request
}

// NewMockContext creates a new MockContext instance with initialized maps
// for URL values, headers, and path parameters. It uses context.Background()
// as the default context.
func NewMockContext[B any]() *MockContext[B] {
	return &MockContext[B]{
		urlValues:  make(url.Values),
		headers:    make(http.Header),
		pathParams: make(map[string]string),
		ctx:        context.Background(),
	}
}

// Body returns the previously set body value. This method always returns
// nil as the error value, as the mock context doesn't perform actual
// deserialization.
func (m *MockContext[B]) Body() (B, error) {
	return m.body, nil
}

// SetBody stores the provided body value for later retrieval via Body().
// This is typically used in tests to simulate request bodies.
func (m *MockContext[B]) SetBody(body B) {
	m.body = body
}

// URLValues returns the mock URL values
func (m *MockContext[B]) URLValues() url.Values {
	return m.urlValues
}

// SetURLValues sets the mock URL values
func (m *MockContext[B]) SetURLValues(values url.Values) {
	m.urlValues = values
}

// Header returns the mock headers
func (m *MockContext[B]) Header() http.Header {
	return m.headers
}

// SetHeader sets a mock header
func (m *MockContext[B]) SetHeader(key, value string) {
	m.headers.Set(key, value)
}

// PathParam returns a mock path parameter
func (m *MockContext[B]) PathParam(name string) string {
	return m.pathParams[name]
}

// SetPathParam sets a mock path parameter
func (m *MockContext[B]) SetPathParam(name, value string) {
	m.pathParams[name] = value
}

// Context returns the mock context
func (m *MockContext[B]) Context() context.Context {
	return m.ctx
}

// SetContext sets the mock context
func (m *MockContext[B]) SetContext(ctx context.Context) {
	m.ctx = ctx
}

// Response returns the mock response writer
func (m *MockContext[B]) Response() http.ResponseWriter {
	return m.response
}

// SetResponse sets the mock response writer
func (m *MockContext[B]) SetResponse(w http.ResponseWriter) {
	m.response = w
}

// Request returns the mock request
func (m *MockContext[B]) Request() *http.Request {
	return m.request
}

// SetRequest sets the mock request
func (m *MockContext[B]) SetRequest(r *http.Request) {
	m.request = r
}
