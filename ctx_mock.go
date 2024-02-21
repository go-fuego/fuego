package fuego

import (
	"fmt"
	"html/template"
	"io/fs"
	"net/http"
	"strconv"
	"strings"
)

type MockContextNoBody = MockContext[any]

// MockContext is used in tests, when the user does not want to use a real http request and response.
type MockContext[B any] struct {
	BaseContextWithBody[B] // Inherits all the methods from BaseContextWithBody

	Req *http.Request
	Res http.ResponseWriter

	fs        fs.FS
	templates *template.Template

	MockBody        B
	MockBodyError   error
	MockQueryParams map[string]string
	MockPathParams  map[string]string
}

var (
	_ Ctx[any]    = &MockContext[any]{}
	_ Ctx[string] = &MockContext[string]{}
	_ Ctx[any]    = &MockContextNoBody{}
)

func (c MockContext[B]) Body() (B, error) {
	return c.MockBody, c.MockBodyError
}

// Render renders the given templates with the given data.
// It returns just an empty string, because the response is written directly to the http.ResponseWriter.
//
// Init templates if not already done.
// This have the side effect of making the Render method static, meaning
// that the templates will be parsed only once, removing
// the need to parse the templates on each request but also preventing
// to dynamically use new templates.
func (c MockContext[B]) Render(templateToExecute string, data any, layoutsGlobs ...string) (HTML, error) {
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

	c.Res.Header().Set("Content-Type", "text/html; charset=utf-8")
	err := c.templates.ExecuteTemplate(c.Res, templateToExecute, data)
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
func (c MockContext[B]) PathParam(name string) string {
	return c.MockPathParams[name]
}

// QueryParams returns the query parameters of the request.
func (c MockContext[B]) QueryParams() map[string]string {
	return c.MockQueryParams
}

// QueryParam returns the query parameter with the given name.
func (c MockContext[B]) QueryParam(name string) string {
	return c.MockQueryParams[name]
}

func (c MockContext[B]) QueryParamIntErr(name string) (int, error) {
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

func (c MockContext[B]) QueryParamInt(name string, defaultValue int) int {
	param, err := c.QueryParamIntErr(name)
	if err != nil {
		return defaultValue
	}

	return param
}

// QueryParamBool returns the query parameter with the given name as a bool.
// If the query parameter does not exist or is not a bool, it returns nil.
// Accepted values are defined as [strconv.ParseBool]
func (c MockContext[B]) QueryParamBoolErr(name string) (bool, error) {
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

func (c MockContext[B]) QueryParamBool(name string, defaultValue bool) bool {
	param, err := c.QueryParamBoolErr(name)
	if err != nil {
		return defaultValue
	}

	return param
}

func (c MockContext[B]) MainLang() string {
	return strings.Split(c.MainLocale(), "-")[0]
}

func (c MockContext[B]) MainLocale() string {
	return strings.Split(c.Req.Header.Get("Accept-Language"), ",")[0]
}

// Request returns the http request.
func (c MockContext[B]) Request() *http.Request {
	return c.Req
}

// Response returns the http response writer.
func (c MockContext[B]) Response() http.ResponseWriter {
	return c.Res
}
