package fuego

import (
	"context"
	"fmt"
	"html/template"
	"io/fs"
	"log/slog"
	"net/http"
	"strings"
)

const (
	maxBodySize = 1048576
)

// Ctx is the context of the request.
// It contains the request body, the path parameters, the query parameters, and the http request.
type Ctx[B any] interface {
	// Body returns the body of the request.
	// If (*B) implements [InTransformer], it will be transformed after deserialization.
	// It caches the result, so it can be called multiple times.
	Body() (B, error)

	// MustBody works like Body, but panics if there is an error.
	MustBody() B

	// PathParam returns the path parameter with the given name.
	// If it does not exist, it returns an empty string.
	// Example:
	//   fuego.Get(s, "/recipes/{recipe_id}", func(c fuego.Ctx[any]) (any, error) {
	//	 	id := c.PathParam("recipe_id")
	//   	...
	//   })
	PathParam(name string) string
	PathParams() map[string]string
	QueryParam(name string) string
	QueryParams() map[string]string

	// Render renders the given templates with the given data.
	// Example:
	//   fuego.Get(s, "/recipes", func(c fuego.Ctx[any]) (any, error) {
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

	Request() *http.Request        // Request returns the underlying http request.
	Response() http.ResponseWriter // Response returns the underlying http response writer.

	// Context returns the context of the request.
	// Same as c.Request().Context().
	// This is the context related to the request, not the context of the server.
	Context() context.Context

	// Redirect redirects to the given url with the given status code.
	// Example:
	//   fuego.Get(s, "/recipes", func(c fuego.Ctx[any]) (any, error) {
	//   	...
	//   	return c.Redirect(301, "/recipes-list")
	//   })
	Redirect(code int, url string) (any, error)
}

func NewContext[B any](w http.ResponseWriter, r *http.Request, options readOptions) *Context[B] {
	c := &Context[B]{
		response: w,
		request:  r,
		readOptions: readOptions{
			DisallowUnknownFields: options.DisallowUnknownFields,
			MaxBodySize:           options.MaxBodySize,
		},
	}

	return c
}

// Context for the request. BodyType is the type of the request body. Please do not use a pointer type as parameter.
type Context[BodyType any] struct {
	body       *BodyType
	request    *http.Request
	response   http.ResponseWriter
	pathParams map[string]string

	fs              fs.FS
	templates       *template.Template
	templatesParsed bool

	readOptions readOptions
}

// SafeShallowCopy returns a safe shallow copy of the context.
// It allows to modify the base context while modifying the request context.
// It is data-safe, meaning that any sensitive data will not be shared between the original context and the copy.
func (c *Context[B]) SafeShallowCopy() *Context[B] {
	c.pathParams = nil
	c.body = nil
	c.request = nil
	c.response = nil

	return c
}

// readOptions are options for reading the request body.
type readOptions struct {
	DisallowUnknownFields bool
	MaxBodySize           int64
	LogBody               bool
}

var _ Ctx[any] = &Context[any]{} // Check that Context implements Ctx.

// Context returns the context of the request.
// Same as c.Request().Context().
func (c Context[B]) Context() context.Context {
	return c.request.Context()
}

func (c Context[B]) Redirect(code int, url string) (any, error) {
	http.Redirect(c.response, c.request, url, code)

	return nil, nil
}

// Render renders the given templates with the given data.
// It returns just an empty string, because the response is written directly to the http.ResponseWriter.
//
// Init templates if not already done.
// This have the side effect of making the Render method static, meaning
// that the templates will be parsed only once, removing
// the need to parse the templates on each request but also preventing
// to dynamically use new templates.
func (c *Context[B]) Render(templateToExecute string, data any, layoutsGlobs ...string) (HTML, error) {
	if !c.templatesParsed {
		layoutsGlobs = append(layoutsGlobs, templateToExecute) // To override all blocks defined in the main template
		cloned := template.Must(c.templates.Clone())
		tmpl, err := cloned.ParseFS(c.fs, layoutsGlobs...)
		if err != nil {
			return "", ErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    fmt.Errorf("error parsing template '%s': %w", layoutsGlobs, err).Error(),
				MoreInfo: map[string]any{
					"templates": layoutsGlobs,
					"help":      "Check that the template exists and have the correct extension.",
				},
			}
		}
		c.templates = template.Must(tmpl.Clone())
		c.templatesParsed = true
	}

	// Get only last template name (for example, with partials/nav/main/nav.partial.html, get nav.partial.html)
	myTemplate := strings.Split(templateToExecute, "/")
	templateToExecute = myTemplate[len(myTemplate)-1]

	c.response.Header().Set("Content-Type", "text/html; charset=utf-8")
	err := c.templates.ExecuteTemplate(c.response, templateToExecute, data)
	if err != nil {
		return "", ErrorResponse{
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
func (c Context[B]) PathParam(name string) string {
	param := c.pathParams[name]
	if param == "" {
		slog.Error("Path parameter might be invalid", "name", name, "valid parameters", c.pathParams)
	}
	return param // TODO go1.22: get (*http.Request) PathValue(name)
}

// PathParams returns the path parameters of the request.
func (c Context[B]) PathParams() map[string]string {
	return nil
}

// QueryParams returns the query parameters of the request.
func (c Context[B]) QueryParams() map[string]string {
	queryParams := c.request.URL.Query()
	params := make(map[string]string)
	for k, v := range queryParams {
		params[k] = v[0]
	}
	return params
}

// QueryParam returns the query parameter with the given name.
func (c Context[B]) QueryParam(name string) string {
	return c.request.URL.Query().Get(name)
}

// Request returns the http request.
func (c Context[B]) Request() *http.Request {
	return c.request
}

// Response returns the http response writer.
func (c Context[B]) Response() http.ResponseWriter {
	return c.response
}

// MustBody works like Body, but panics if there is an error.
func (c *Context[B]) MustBody() B {
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
func (c *Context[B]) Body() (B, error) {
	if c.body != nil {
		return *c.body, nil
	}

	// Limit the size of the request body.
	if c.readOptions.MaxBodySize != 0 {
		c.request.Body = http.MaxBytesReader(nil, c.request.Body, c.readOptions.MaxBodySize)
	}

	var body B
	var err error
	switch c.request.Header.Get("Content-Type") {
	case "text/plain":
		s, errReadingString := readString[string](c.request.Body, c.readOptions)
		body = any(s).(B)
		err = errReadingString
	case "application/x-www-form-urlencoded", "multipart/form-data":
		body, err = readURLEncoded[B](c.request, c.readOptions)
	case "application/json":
		fallthrough
	default:
		body, err = readJSON[B](c.request.Body, c.readOptions)
	}

	c.body = &body

	return body, err
}
