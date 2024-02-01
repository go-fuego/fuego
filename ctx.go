package fuego

import (
	"context"
	"fmt"
	"html/template"
	"io/fs"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const (
	maxBodySize = 1048576
)

// ctx is the context of the request.
// It contains the request body, the path parameters, the query parameters, and the http request.
// Please do not use a pointer type as parameter.
type ctx[B any] interface {
	// Body returns the body of the request.
	// If (*B) implements [InTransformer], it will be transformed after deserialization.
	// It caches the result, so it can be called multiple times.
	Body() (B, error)

	// MustBody works like Body, but panics if there is an error.
	MustBody() B

	// PathParam returns the path parameter with the given name.
	// If it does not exist, it returns an empty string.
	// Example:
	//   fuego.Get(s, "/recipes/{recipe_id}", func(c fuego.ContextNoBody) (any, error) {
	//	 	id := c.PathParam("recipe_id")
	//   	...
	//   })
	PathParam(name string) string
	PathParams() map[string]string
	QueryParam(name string) string
	QueryParamInt(name string, defaultValue int) int // If the query parameter does not exist or is not an int, it returns the default given value. Use [Ctx.QueryParamIntErr] if you want to know if the query parameter is erroneous.
	QueryParamIntErr(name string) (int, error)
	QueryParamBool(name string, defaultValue bool) bool // If the query parameter does not exist or is not a bool, it returns the default given value. Use [Ctx.QueryParamBoolErr] if you want to know if the query parameter is erroneous.
	QueryParamBoolErr(name string) (bool, error)
	QueryParams() map[string]string

	MainLang() string   // ex: fr. MainLang returns the main language of the request. It is the first language of the Accept-Language header. To get the main locale (ex: fr-CA), use [Ctx.MainLocale].
	MainLocale() string // ex: en-US. MainLocale returns the main locale of the request. It is the first locale of the Accept-Language header. To get the main language (ex: en), use [Ctx.MainLang].

	// Render renders the given templates with the given data.
	// Example:
	//   fuego.Get(s, "/recipes", func(c fuego.ContextNoBody) (any, error) {
	//   	recipes, _ := rs.Queries.GetRecipes(c.Context())
	//   		...
	//   	return c.Render("pages/recipes.page.html", recipes)
	//   })
	// For the Go templates reference, see https://pkg.go.dev/html/template
	//
	// [templateGlobsToOverride] is a list of templates to override.
	// For example, if you have 2 conflicting templates
	//   - with the same name "partials/aaa/nav.partial.html" and "partials/bbb/nav.partial.html"
	//   - or two templates with different names but that define the same block "page" for example,
	// and you want to override one above the other, you can do:
	//   c.Render("admin.page.html", recipes, "partials/aaa/nav.partial.html")
	// By default, [templateToExecute] is added to the list of templates to override.
	Render(templateToExecute string, data any, templateGlobsToOverride ...string) (HTML, error)

	Context() context.Context

	Request() *http.Request        // Request returns the underlying http request.
	Response() http.ResponseWriter // Response returns the underlying http response writer.

	// SetStatus sets the status code of the response.
	// Alias to http.ResponseWriter.WriteHeader.
	SetStatus(code int)

	// Redirect redirects to the given url with the given status code.
	// Example:
	//   fuego.Get(s, "/recipes", func(c fuego.ContextNoBody) (any, error) {
	//   	...
	//   	return c.Redirect(301, "/recipes-list")
	//   })
	Redirect(code int, url string) (any, error)
}

// NewContext returns a new context. It is used internally by Fuego. You probably want to use Ctx[B] instead.
func NewContext[B any](w http.ResponseWriter, r *http.Request, options readOptions) *ContextWithBody[B] {
	c := &ContextWithBody[B]{
		ContextNoBody: ContextNoBody{
			response: w,
			request:  r,
			readOptions: readOptions{
				DisallowUnknownFields: options.DisallowUnknownFields,
				MaxBodySize:           options.MaxBodySize,
			},
		},
	}

	return c
}

// ContextWithBody is used internally by Fuego. You probably want to use Ctx[B] instead. Please do not use a pointer type as parameter.
type ContextWithBody[Body any] struct {
	body *Body // Cache the body in request context, because it is not possible to read an http request body multiple times.
	ContextNoBody
}

// ContextNoBody is used when the controller does not have a body.
// It used as a base context for other Context types.
type ContextNoBody struct {
	request    *http.Request
	response   http.ResponseWriter
	pathParams map[string]string

	fs        fs.FS
	templates *template.Template

	readOptions readOptions
}

func (c ContextNoBody) Body() (any, error) {
	slog.Warn("this method should not be called. It probably happened because you passed the context to another controller.")
	return body[map[string]any](c)
}

func (c ContextNoBody) MustBody() any {
	b, err := c.Body()
	if err != nil {
		panic(err)
	}
	return b
}

// SetStatus sets the status code of the response.
// Alias to http.ResponseWriter.WriteHeader.
func (c ContextNoBody) SetStatus(code int) {
	c.response.WriteHeader(code)
}

// readOptions are options for reading the request body.
type readOptions struct {
	DisallowUnknownFields bool
	MaxBodySize           int64
	LogBody               bool
}

var (
	_ ctx[any]    = &ContextWithBody[any]{}    // Check that ContextWithBody[any] implements Ctx.
	_ ctx[string] = &ContextWithBody[string]{} // Check that ContextWithBody[string] implements Ctx.
	_ ctx[any]    = &ContextNoBody{}           // Check that ContextNoBody implements Ctx.
	_ ctx[any]    = ContextNoBody{}            // Check that ContextNoBody implements Ctx.
)

func (c ContextNoBody) Redirect(code int, url string) (any, error) {
	http.Redirect(c.response, c.request, url, code)

	return nil, nil
}

func (c ContextNoBody) Context() context.Context {
	return c.request.Context()
}

// Render renders the given templates with the given data.
// It returns just an empty string, because the response is written directly to the http.ResponseWriter.
//
// Init templates if not already done.
// This have the side effect of making the Render method static, meaning
// that the templates will be parsed only once, removing
// the need to parse the templates on each request but also preventing
// to dynamically use new templates.
func (c ContextNoBody) Render(templateToExecute string, data any, layoutsGlobs ...string) (HTML, error) {
	if strings.Contains(templateToExecute, "/") || strings.Contains(templateToExecute, "*") {

		layoutsGlobs = append(layoutsGlobs, templateToExecute) // To override all blocks defined in the main template
		cloned := template.Must(c.templates.Clone())
		tmpl, err := cloned.ParseFS(c.fs, layoutsGlobs...)
		if err != nil {
			return "", HTTPError{
				StatusCode: http.StatusInternalServerError,
				Message:    fmt.Errorf("error parsing template '%s': %w", layoutsGlobs, err).Error(),
				MoreInfo: map[string]any{
					"templates": layoutsGlobs,
					"help":      "Check that the template exists and have the correct extension.",
				},
			}
		}
		c.templates = template.Must(tmpl.Clone())
	}

	// Get only last template name (for example, with partials/nav/main/nav.partial.html, get nav.partial.html)
	myTemplate := strings.Split(templateToExecute, "/")
	templateToExecute = myTemplate[len(myTemplate)-1]

	c.response.Header().Set("Content-Type", "text/html; charset=utf-8")
	err := c.templates.ExecuteTemplate(c.response, templateToExecute, data)
	if err != nil {
		return "", HTTPError{
			StatusCode: http.StatusInternalServerError,
			Message:    fmt.Errorf("error executing template '%s': %w", templateToExecute, err).Error(),
			MoreInfo: map[string]any{
				"templates": layoutsGlobs,
				"help":      "Check that the template exists and have the correct extension.",
			},
		}
	}

	return "", err
}

// PathParams returns the path parameters of the request.
func (c ContextNoBody) PathParam(name string) string {
	param := c.pathParams[name]
	if param == "" {
		slog.Error("Path parameter might be invalid", "name", name, "valid parameters", c.pathParams)
	}
	return param // TODO go1.22: get (*http.Request) PathValue(name)
}

// PathParams returns the path parameters of the request.
func (c ContextNoBody) PathParams() map[string]string {
	return nil
}

type QueryParamNotFoundError struct {
	ParamName string
}

func (e QueryParamNotFoundError) Error() string {
	return fmt.Errorf("param %s not found", e.ParamName).Error()
}

type QueryParamInvalidTypeError struct {
	ParamName    string
	ParamValue   string
	ExpectedType string
	Err          error
}

func (e QueryParamInvalidTypeError) Error() string {
	return fmt.Errorf("param %s=%s is not of type %s: %w", e.ParamName, e.ParamValue, e.ExpectedType, e.Err).Error()
}

// QueryParams returns the query parameters of the request.
func (c ContextNoBody) QueryParams() map[string]string {
	queryParams := c.request.URL.Query()
	params := make(map[string]string)
	for k, v := range queryParams {
		params[k] = v[0]
	}
	return params
}

// QueryParam returns the query parameter with the given name.
func (c ContextNoBody) QueryParam(name string) string {
	return c.request.URL.Query().Get(name)
}

func (c ContextNoBody) QueryParamIntErr(name string) (int, error) {
	param := c.QueryParam(name)
	if param == "" {
		return 0, QueryParamNotFoundError{ParamName: name}
	}

	i, err := strconv.Atoi(param)
	if err != nil {
		return 0, QueryParamInvalidTypeError{
			ParamName:    name,
			ParamValue:   param,
			ExpectedType: "int",
			Err:          err,
		}
	}

	return i, nil
}

func (c ContextNoBody) QueryParamInt(name string, defaultValue int) int {
	param, err := c.QueryParamIntErr(name)
	if err != nil {
		return defaultValue
	}

	return param
}

// QueryParamBool returns the query parameter with the given name as a bool.
// If the query parameter does not exist or is not a bool, it returns nil.
// Accepted values are defined as [strconv.ParseBool]
func (c ContextNoBody) QueryParamBoolErr(name string) (bool, error) {
	param := c.QueryParam(name)
	if param == "" {
		return false, QueryParamNotFoundError{ParamName: name}
	}

	b, err := strconv.ParseBool(param)
	if err != nil {
		return false, QueryParamInvalidTypeError{
			ParamName:    name,
			ParamValue:   param,
			ExpectedType: "bool",
			Err:          err,
		}
	}
	return b, nil
}

func (c ContextNoBody) QueryParamBool(name string, defaultValue bool) bool {
	param, err := c.QueryParamBoolErr(name)
	if err != nil {
		return defaultValue
	}

	return param
}

func (c ContextNoBody) MainLang() string {
	return strings.Split(c.MainLocale(), "-")[0]
}

func (c ContextNoBody) MainLocale() string {
	return strings.Split(c.request.Header.Get("Accept-Language"), ",")[0]
}

// Request returns the http request.
func (c ContextNoBody) Request() *http.Request {
	return c.request
}

// Response returns the http response writer.
func (c ContextNoBody) Response() http.ResponseWriter {
	return c.response
}

// MustBody works like Body, but panics if there is an error.
func (c *ContextWithBody[B]) MustBody() B {
	b, err := c.Body()
	if err != nil {
		panic(err)
	}
	return b
}

// Body returns the body of the request.
// If (*B) implements [InTransformer], it will be transformed after deserialization.
// It caches the result, so it can be called multiple times.
// The reason why the body is cached is because it is not possible to read an http request body multiple times, not because of performance.
// For decoding, it uses the Content-Type header. If it is not set, defaults to application/json.
func (c *ContextWithBody[B]) Body() (B, error) {
	if c.body != nil {
		return *c.body, nil
	}

	body, err := body[B](c.ContextNoBody)
	c.body = &body
	return body, err
}

func body[B any](c ContextNoBody) (B, error) {
	// Limit the size of the request body.
	if c.readOptions.MaxBodySize != 0 {
		c.request.Body = http.MaxBytesReader(nil, c.request.Body, c.readOptions.MaxBodySize)
	}

	timeDeserialize := time.Now()

	var body B
	var err error
	switch c.request.Header.Get("Content-Type") {
	case "text/plain":
		s, errReadingString := readString[string](c.request.Context(), c.request.Body, c.readOptions)
		body = any(s).(B)
		err = errReadingString
	case "application/x-www-form-urlencoded", "multipart/form-data":
		body, err = readURLEncoded[B](c.request, c.readOptions)
	case "application/json":
		fallthrough
	default:
		body, err = readJSON[B](c.request.Context(), c.request.Body, c.readOptions)
	}

	c.response.Header().Add("Server-Timing", Timing{"deserialize", time.Since(timeDeserialize), "controller > deserialize"}.String())

	return body, err
}
