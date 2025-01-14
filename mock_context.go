package fuego

import (
	"context"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
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
	cookies    map[string]*http.Cookie
	params     map[string]OpenAPIParam
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
		cookies:    make(map[string]*http.Cookie),
		params:     make(map[string]OpenAPIParam),
	}
}

// GetOpenAPIParams returns the OpenAPI parameters for validation
func (m *MockContext[B]) GetOpenAPIParams() map[string]OpenAPIParam {
	return m.params
}

// SetOpenAPIParam sets an OpenAPI parameter for validation
func (m *MockContext[B]) SetOpenAPIParam(name string, param OpenAPIParam) {
	m.params[name] = param
}

// HasQueryParam checks if a query parameter exists
func (m *MockContext[B]) HasQueryParam(key string) bool {
	_, exists := m.urlValues[key]
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
	return m.body, nil
}

// MustBody returns the body or panics if there's an error
func (m *MockContext[B]) MustBody() B {
	return m.body
}

// SetBody stores the provided body value for later retrieval
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

// Header returns the value of the specified header
func (m *MockContext[B]) Header(key string) string {
	return m.headers.Get(key)
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

// Deadline implements context.Context
func (m *MockContext[B]) Deadline() (deadline time.Time, ok bool) {
	return m.ctx.Deadline()
}

// Done implements context.Context
func (m *MockContext[B]) Done() <-chan struct{} {
	return m.ctx.Done()
}

// Err implements context.Context
func (m *MockContext[B]) Err() error {
	return m.ctx.Err()
}

// Value implements context.Context
func (m *MockContext[B]) Value(key any) any {
	return m.ctx.Value(key)
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

// QueryParam returns the value of the specified query parameter
func (m *MockContext[B]) QueryParam(name string) string {
	return m.urlValues.Get(name)
}

// QueryParamArr returns the values of the specified query parameter
func (m *MockContext[B]) QueryParamArr(name string) []string {
	return m.urlValues[name]
}

// QueryParamInt returns the value of the specified query parameter as an integer
func (m *MockContext[B]) QueryParamInt(name string) int {
	val := m.QueryParam(name)
	if val == "" {
		return 0
	}
	i, _ := strconv.Atoi(val)
	return i
}

// QueryParamIntErr returns the value of the specified query parameter as an integer and any error
func (m *MockContext[B]) QueryParamIntErr(name string) (int, error) {
	val := m.QueryParam(name)
	if val == "" {
		return 0, nil
	}
	return strconv.Atoi(val)
}

// QueryParamBool returns the value of the specified query parameter as a boolean
func (m *MockContext[B]) QueryParamBool(name string) bool {
	val := m.QueryParam(name)
	if val == "" {
		return false
	}
	b, _ := strconv.ParseBool(val)
	return b
}

// QueryParamBoolErr returns the value of the specified query parameter as a boolean and any error
func (m *MockContext[B]) QueryParamBoolErr(name string) (bool, error) {
	val := m.QueryParam(name)
	if val == "" {
		return false, nil
	}
	return strconv.ParseBool(val)
}

// QueryParams returns all query parameters
func (m *MockContext[B]) QueryParams() url.Values {
	return m.urlValues
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
	return nil, nil
}
