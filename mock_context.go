package fuego

import (
	"context"
	"net/http"
	"net/url"
	"strings"

	"github.com/go-fuego/fuego/internal"
)

// MockContext provides a framework-agnostic implementation of ContextWithBody
// for testing purposes. It allows testing controllers without depending on
// specific web frameworks like Gin or Echo.
type MockContext[B any] struct {
	internal.CommonContext[B]
	BodyData       B // Public for easy testing
	HeadersData    http.Header
	PathParamsData map[string]string
	ResponseData   http.ResponseWriter
	RequestData    *http.Request
	CookiesData    map[string]*http.Cookie
}

// NewMockContext creates a new MockContext instance with initialized maps
// for URL values, headers, and path parameters. It uses context.Background()
// as the default context.
func NewMockContext[B any]() *MockContext[B] {
	return &MockContext[B]{
		CommonContext: internal.CommonContext[B]{
			CommonCtx:     context.Background(),
			UrlValues:     make(url.Values),
			OpenAPIParams: make(map[string]internal.OpenAPIParam),
		},
		HeadersData:    make(http.Header),
		PathParamsData: make(map[string]string),
		CookiesData:    make(map[string]*http.Cookie),
	}
}

// Body returns the body value - implements ContextWithBody
func (m *MockContext[B]) Body() (B, error) {
	return m.BodyData, nil
}

// MustBody returns the body or panics if there's an error - implements ContextWithBody
func (m *MockContext[B]) MustBody() B {
	return m.BodyData
}

// Header returns the value of the specified header - implements ContextWithBody
func (m *MockContext[B]) Header(key string) string {
	return m.HeadersData.Get(key)
}

// HasHeader checks if a header exists - implements ContextWithBody
func (m *MockContext[B]) HasHeader(key string) bool {
	_, exists := m.HeadersData[key]
	return exists
}

// HasCookie checks if a cookie exists - implements ContextWithBody
func (m *MockContext[B]) HasCookie(key string) bool {
	_, exists := m.CookiesData[key]
	return exists
}

// PathParam returns a mock path parameter - implements ContextWithBody
func (m *MockContext[B]) PathParam(name string) string {
	return m.PathParamsData[name]
}

// Response returns the mock response writer - implements ContextWithBody
func (m *MockContext[B]) Response() http.ResponseWriter {
	return m.ResponseData
}

// Request returns the mock request - implements ContextWithBody
func (m *MockContext[B]) Request() *http.Request {
	return m.RequestData
}

// Cookie returns a mock cookie - implements ContextWithBody
func (m *MockContext[B]) Cookie(name string) (*http.Cookie, error) {
	cookie, exists := m.CookiesData[name]
	if !exists {
		return nil, http.ErrNoCookie
	}
	return cookie, nil
}

// MainLang returns the main language from Accept-Language header - implements ContextWithBody
func (m *MockContext[B]) MainLang() string {
	lang := m.HeadersData.Get("Accept-Language")
	if lang == "" {
		return ""
	}
	return strings.Split(strings.Split(lang, ",")[0], "-")[0]
}

// MainLocale returns the main locale from Accept-Language header - implements ContextWithBody
func (m *MockContext[B]) MainLocale() string {
	return m.HeadersData.Get("Accept-Language")
}

// Redirect returns a redirect response - implements ContextWithBody
func (m *MockContext[B]) Redirect(code int, url string) (any, error) {
	if m.ResponseData != nil {
		http.Redirect(m.ResponseData, m.RequestData, url, code)
	}
	return nil, nil
}

// Render is a mock implementation that does nothing - implements ContextWithBody
func (m *MockContext[B]) Render(templateToExecute string, data any, templateGlobsToOverride ...string) (CtxRenderer, error) {
	return nil, nil
}

// SetStatus sets the response status code - implements ContextWithBody
func (m *MockContext[B]) SetStatus(code int) {
	if m.ResponseData != nil {
		m.ResponseData.WriteHeader(code)
	}
}

// SetCookie sets a mock cookie - implements ContextWithBody
func (m *MockContext[B]) SetCookie(cookie http.Cookie) {
	m.CookiesData[cookie.Name] = &cookie
}

// SetHeader sets a mock header - implements ContextWithBody
func (m *MockContext[B]) SetHeader(key, value string) {
	m.HeadersData.Set(key, value)
}
