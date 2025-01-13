package fuego

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/go-fuego/fuego/internal"
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

// Header returns the value of the header with the given key
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

// Cookie returns a cookie by name
func (m *MockContext[B]) Cookie(name string) (*http.Cookie, error) {
	if cookie, exists := m.cookies[name]; exists {
		return cookie, nil
	}
	return nil, http.ErrNoCookie
}

// SetCookie sets a cookie for testing
func (m *MockContext[B]) SetCookie(cookie http.Cookie) {
	m.cookies[cookie.Name] = &cookie
}

// Deadline returns the time when work done on behalf of this context
// should be canceled. In this mock implementation, we return no deadline.
func (m *MockContext[B]) Deadline() (deadline time.Time, ok bool) {
	return time.Time{}, false
}

// Done returns a channel that's closed when work done on behalf of this
// context should be canceled. In this mock implementation, we return nil
// which means the context can never be canceled.
func (m *MockContext[B]) Done() <-chan struct{} {
	return nil
}

// Err returns nil since this mock context never returns errors
func (m *MockContext[B]) Err() error {
	return nil
}

// EmptyBody represents an empty request body
type EmptyBody struct{}

// GetOpenAPIParams returns an empty map since this is just a mock
func (m *MockContext[B]) GetOpenAPIParams() map[string]internal.OpenAPIParam {
	return make(map[string]internal.OpenAPIParam)
}

// HasCookie checks if a cookie with the given name exists
func (m *MockContext[B]) HasCookie(name string) bool {
	_, exists := m.cookies[name]
	return exists
}

// HasHeader checks if a header with the given key exists
func (m *MockContext[B]) HasHeader(key string) bool {
	_, exists := m.headers[key]
	return exists
}

// HasQueryParam checks if a query parameter with the given key exists
func (m *MockContext[B]) HasQueryParam(key string) bool {
	_, exists := m.urlValues[key]
	return exists
}

// MainLang returns the main language for the request (e.g., "en").
// In this mock implementation, we'll return "en" as default.
func (m *MockContext[B]) MainLang() string {
	// Get language from Accept-Language header or return default
	if lang := m.headers.Get("Accept-Language"); lang != "" {
		return lang[:2] // Take first two chars for language code
	}
	return "en" // Default to English
}

// MainLocale returns the main locale for the request (e.g., "en-US").
// In this mock implementation, we'll return "en-US" as default.
func (m *MockContext[B]) MainLocale() string {
	// Get locale from Accept-Language header or return default
	if locale := m.headers.Get("Accept-Language"); locale != "" {
		return locale
	}
	return "en-US" // Default to English (US)
}

// MustBody returns the body directly, without error handling.
// In this mock implementation, we simply return the body since we know it's valid.
func (m *MockContext[B]) MustBody() B {
	return m.body
}

// QueryParam returns the value of the query parameter with the given key.
// If there are multiple values, it returns the first one.
// If the parameter doesn't exist, it returns an empty string.
func (m *MockContext[B]) QueryParam(key string) string {
	return m.urlValues.Get(key)
}

// QueryParamArr returns all values for the query parameter with the given key.
// If the parameter doesn't exist, it returns an empty slice.
func (m *MockContext[B]) QueryParamArr(key string) []string {
	return m.urlValues[key]
}

// QueryParamBool returns the boolean value of the query parameter with the given key.
// Returns true for "1", "t", "T", "true", "TRUE", "True"
// Returns false for "0", "f", "F", "false", "FALSE", "False"
// Returns false for any other value
func (m *MockContext[B]) QueryParamBool(key string) bool {
	v := m.urlValues.Get(key)
	switch v {
	case "1", "t", "T", "true", "TRUE", "True":
		return true
	default:
		return false
	}
}

// QueryParamBoolErr returns the boolean value of the query parameter with the given key
// and an error if the value is not a valid boolean.
func (m *MockContext[B]) QueryParamBoolErr(key string) (bool, error) {
	v := m.urlValues.Get(key)
	switch v {
	case "1", "t", "T", "true", "TRUE", "True":
		return true, nil
	case "0", "f", "F", "false", "FALSE", "False":
		return false, nil
	case "":
		return false, nil // Parameter not found
	default:
		return false, fmt.Errorf("invalid boolean value: %s", v)
	}
}

// QueryParamInt returns the integer value of the query parameter with the given key.
// Returns 0 if the parameter doesn't exist or cannot be parsed as an integer.
func (m *MockContext[B]) QueryParamInt(key string) int {
	v := m.urlValues.Get(key)
	if v == "" {
		return 0
	}
	i, err := strconv.Atoi(v)
	if err != nil {
		return 0
	}
	return i
}

// QueryParamIntErr returns the integer value of the query parameter with the given key
// and an error if the value cannot be parsed as an integer.
func (m *MockContext[B]) QueryParamIntErr(key string) (int, error) {
	v := m.urlValues.Get(key)
	if v == "" {
		return 0, nil // Parameter not found
	}
	i, err := strconv.Atoi(v)
	if err != nil {
		return 0, fmt.Errorf("invalid integer value: %s", v)
	}
	return i, nil
}

// QueryParams returns all query parameters.
// This is an alias for URLValues() for compatibility with the interface.
func (m *MockContext[B]) QueryParams() url.Values {
	return m.urlValues
}

// Redirect performs a redirect to the specified URL with the given status code.
// In this mock implementation, we store the redirect information for testing.
func (m *MockContext[B]) Redirect(code int, url string) (any, error) {
	if m.response != nil {
		m.response.Header().Set("Location", url)
		m.response.WriteHeader(code)
	}
	return nil, nil
}

// mockRenderer implements CtxRenderer for testing
type mockRenderer struct {
	data     any
	template string
	layouts  []string
}

// Render implements the CtxRenderer interface
func (r *mockRenderer) Render(ctx context.Context, w io.Writer) error {
	// In a real implementation, this would render the template
	// For testing, we just write a success status
	if hw, ok := w.(http.ResponseWriter); ok {
		hw.WriteHeader(http.StatusOK)
	}
	return nil
}

// Render renders the template with the given name and data.
// In this mock implementation, we just store the data for testing.
func (m *MockContext[B]) Render(templateName string, data any, layouts ...string) (CtxRenderer, error) {
	if m.response != nil {
		// In a real implementation, this would render the template with layouts
		// For testing, we just store the data
	}
	return &mockRenderer{
		data:     data,
		template: templateName,
		layouts:  layouts,
	}, nil
}

// SetStatus sets the HTTP status code for the response.
// In this mock implementation, we set the status code if a response writer is available.
func (m *MockContext[B]) SetStatus(code int) {
	if m.response != nil {
		m.response.WriteHeader(code)
	}
}

// Value returns the value associated with this context for key, or nil
// if no value is associated with key. In this mock implementation,
// we delegate to the underlying context.
func (m *MockContext[B]) Value(key any) any {
	return m.ctx.Value(key)
}
