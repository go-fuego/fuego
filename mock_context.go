package fuego

import (
	"context"
	"net/http"
	"net/url"
	"strings"

	"github.com/go-fuego/fuego/internal"
)

// NewMockContext creates a new MockContext instance with initialized maps
// for URL values, headers, and path parameters. It uses context.Background()
// as the default context.
func NewMockContext[B any]() *MockContext[B] {
	return &MockContext[B]{
		CommonContext: internal.CommonContext[B]{
			CommonCtx: context.Background(),
		},
		headers:    make(http.Header),
		pathParams: make(map[string]string),
		cookies:    make(map[string]*http.Cookie),
	}
}

// MockContext provides a framework-agnostic implementation of ContextWithBody
// for testing purposes. It allows testing controllers without depending on
// specific web frameworks like Gin or Echo.
type MockContext[B any] struct {
	internal.CommonContext[B]

	body B

	headers    http.Header
	pathParams map[string]string
	response   http.ResponseWriter
	request    *http.Request
	cookies    map[string]*http.Cookie
}

var _ ContextWithBody[string] = &MockContext[string]{}

// SetOpenAPIParam sets an OpenAPI parameter for validation
func (m *MockContext[B]) SetOpenAPIParam(name string, param OpenAPIParam) {
	m.CommonContext.OpenAPIParams[name] = param
}

// HasQueryParam checks if a query parameter exists
func (m *MockContext[B]) HasQueryParam(key string) bool {
	_, exists := m.UrlValues[key]
	return exists
}

// HasHeader checks if a header exists
func (m *MockContext[B]) HasHeader(key string) bool {
	_, exists := m.headers[key]
	return exists
}

// HasCookie checks if a cookie exists
func (m *MockContext[B]) HasCookie(key string) bool {
	_, exists := m.cookies[key]
	return exists
}

// Body returns the previously set body value
func (m *MockContext[B]) Body() (B, error) {
	return TransformAndValidate(m, m.body)
}

// MustBody returns the body or panics if there's an error
func (m *MockContext[B]) MustBody() B {
	return m.body
}

// SetBody stores the provided body value for later retrieval
func (m *MockContext[B]) SetBody(body B) {
	m.body = body
}

// SetQueryParams sets the mock URL values
func (m *MockContext[B]) SetQueryParams(values url.Values) {
	m.UrlValues = values
}

// SetQueryParam sets a single mock URL value
func (m *MockContext[B]) SetQueryParam(key, value string) {
	m.UrlValues.Set(key, value)
}

// Header returns the value of the specified header
func (m *MockContext[B]) Header(key string) string {
	return m.headers.Get(key)
}

// SetHeader sets a mock header
func (m *MockContext[B]) SetHeader(key, value string) {
	m.headers.Set(key, value)
}

// GetHeaders returns all headers (helper method for testing)
func (m *MockContext[B]) GetHeaders() http.Header {
	return m.headers
}

// PathParam returns a mock path parameter
func (m *MockContext[B]) PathParam(name string) string {
	return m.pathParams[name]
}

// SetPathParam sets a mock path parameter
func (m *MockContext[B]) SetPathParam(name, value string) {
	m.pathParams[name] = value
}

// SetContext sets the mock context
func (m *MockContext[B]) SetContext(ctx context.Context) {
	m.CommonContext.CommonCtx = ctx
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

// Cookie returns a mock cookie
func (m *MockContext[B]) Cookie(name string) (*http.Cookie, error) {
	cookie, exists := m.cookies[name]
	if !exists {
		return nil, http.ErrNoCookie
	}
	return cookie, nil
}

// SetCookie sets a mock cookie
func (m *MockContext[B]) SetCookie(cookie http.Cookie) {
	m.cookies[cookie.Name] = &cookie
}

// MainLang returns the main language from Accept-Language header
func (m *MockContext[B]) MainLang() string {
	lang := m.headers.Get("Accept-Language")
	if lang == "" {
		return ""
	}
	return strings.Split(strings.Split(lang, ",")[0], "-")[0]
}

// MainLocale returns the main locale from Accept-Language header
func (m *MockContext[B]) MainLocale() string {
	return m.headers.Get("Accept-Language")
}

// SetStatus sets the response status code
func (m *MockContext[B]) SetStatus(code int) {
	if m.response != nil {
		m.response.WriteHeader(code)
	}
}

// Redirect returns a redirect response
func (m *MockContext[B]) Redirect(code int, url string) (any, error) {
	if m.response != nil {
		http.Redirect(m.response, m.request, url, code)
	}
	return nil, nil
}

// Render is a mock implementation that does nothing
func (m *MockContext[B]) Render(templateToExecute string, data any, templateGlobsToOverride ...string) (CtxRenderer, error) {
	panic("not implemented")
}
