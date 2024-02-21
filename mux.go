package fuego

import (
	"log/slog"
	"net/http"
	"reflect"
	"runtime"
	"slices"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
)

// Group allows to group routes under a common path.
// Middlewares are scoped to the group.
// For example:
//
//	s := fuego.NewServer()
//	viewsRoutes := fuego.Group("")
//	apiRoutes := fuego.Group("/api")
//	// Registering a middlewares scoped to /api only
//	fuego.Use(apiRoutes, myMiddleware)
//	// Registering a route under /api/users
//	fuego.Get(apiRoutes, "/users", func(c fuego.ContextNoBody) (ans, error) {
//		return ans{Ans: "users"}, nil
//	})
//	s.Run()
func Group(s *Server, path string) *Server {
	if path == "/" {
		path = ""
	} else if path != "" && path[len(path)-1] == '/' {
		slog.Warn("Group path should not end with a slash.", "path", path+"/", "new", path)
	}

	ss := *s
	newServer := &ss
	newServer.basePath += path

	return newServer
}

type Route[ResponseBody any, RequestBody any] struct {
	operation *openapi3.Operation
}

// Capture all methods (GET, POST, PUT, PATCH, DELETE) and register a controller.
func All[T any, B any, Contexted ctx[B]](s *Server, path string, controller func(Contexted) (T, error), middlewares ...func(http.Handler) http.Handler) Route[T, B] {
	for _, method := range []string{http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodDelete} {
		Register[T](s, method, path, controller, middlewares...)
	}
	return Register[T](s, http.MethodGet, path, controller, middlewares...)
}

func Get[T any, B any, Contexted ctx[B]](s *Server, path string, controller func(Contexted) (T, error), middlewares ...func(http.Handler) http.Handler) Route[T, B] {
	return Register[T](s, http.MethodGet, path, controller, middlewares...)
}

func Post[T any, B any, Contexted ctx[B]](s *Server, path string, controller func(Contexted) (T, error), middlewares ...func(http.Handler) http.Handler) Route[T, B] {
	return Register[T](s, http.MethodPost, path, controller, middlewares...)
}

func Delete[T any, B any, Contexted ctx[B]](s *Server, path string, controller func(Contexted) (T, error), middlewares ...func(http.Handler) http.Handler) Route[T, B] {
	return Register[T](s, http.MethodDelete, path, controller, middlewares...)
}

func Put[T any, B any, Contexted ctx[B]](s *Server, path string, controller func(Contexted) (T, error), middlewares ...func(http.Handler) http.Handler) Route[T, B] {
	return Register[T](s, http.MethodPut, path, controller, middlewares...)
}

func Patch[T any, B any, Contexted ctx[B]](s *Server, path string, controller func(Contexted) (T, error), middlewares ...func(http.Handler) http.Handler) Route[T, B] {
	return Register[T](s, http.MethodPatch, path, controller, middlewares...)
}

// Registers route into the default mux.
func Register[T any, B any, Contexted ctx[B]](s *Server, method string, path string, controller func(Contexted) (T, error), middlewares ...func(http.Handler) http.Handler) Route[T, B] {
	fullPath := method + " " + s.basePath + path

	slog.Debug("registering openapi controller " + fullPath)

	route := register[T, B](s, method, path, httpHandler[T, B](s, controller), middlewares...)

	name, nameWithPath := funcName(controller)
	route.operation.Summary = name
	route.operation.Description = "controller: " + nameWithPath
	route.operation.OperationID = fullPath + ":" + name
	return route
}

func register[T any, B any](s *Server, method string, path string, controller http.Handler, middlewares ...func(http.Handler) http.Handler) Route[T, B] {
	fullPath := method + " " + s.basePath + path

	allMiddlewares := append(middlewares, s.middlewares...)
	s.Mux.Handle(fullPath, withMiddlewares(controller, allMiddlewares...))

	operation, err := RegisterOpenAPIOperation[T, B](s, method, s.basePath+path)
	if err != nil {
		slog.Warn("error documenting openapi operation", "error", err)
	}

	return Route[T, B]{
		operation: operation,
	}
}

func (r Route[ResponseBody, RequestBody]) WithDescription(description string) Route[ResponseBody, RequestBody] {
	r.operation.Description = description
	return r
}

func (r Route[ResponseBody, RequestBody]) WithSummary(summary string) Route[ResponseBody, RequestBody] {
	r.operation.Summary = summary
	return r
}

func (r Route[ResponseBody, RequestBody]) SetTags(tags ...string) Route[ResponseBody, RequestBody] {
	r.operation.Tags = tags
	return r
}

type OpenAPIParam struct {
	Required bool
	Example  string
	Type     string // "query", "header", "cookie"
}

// Param registers a parameter for the route.
// The paramType can be "query", "header" or "cookie".
// [Cookie], [Header], [QueryParam] are shortcuts for Param.
func (r Route[ResponseBody, RequestBody]) Param(paramType, name, description string, params ...OpenAPIParam) Route[ResponseBody, RequestBody] {
	openapiParam := openapi3.NewHeaderParameter(name)
	openapiParam.Description = description
	openapiParam.Schema = openapi3.NewStringSchema().NewRef()
	openapiParam.In = paramType

	for _, param := range params {
		openapiParam.Required = param.Required
		openapiParam.Example = param.Example
		openapiParam.In = param.Type
	}

	r.operation.AddParameter(openapiParam)

	return r
}

// Header registers a header parameter for the route.
func (r Route[ResponseBody, RequestBody]) Header(name, description string, params ...OpenAPIParam) Route[ResponseBody, RequestBody] {
	r.Param("header", name, description, params...)
	return r
}

// Cookie registers a cookie parameter for the route.
func (r Route[ResponseBody, RequestBody]) Cookie(name, description string, params ...OpenAPIParam) Route[ResponseBody, RequestBody] {
	r.Param("cookie", name, description, params...)
	return r
}

// QueryParam registers a query parameter for the route.
func (r Route[ResponseBody, RequestBody]) QueryParam(name, description string, params ...OpenAPIParam) Route[ResponseBody, RequestBody] {
	r.Param("query", name, description, params...)
	return r
}

func (r Route[ResponseBody, RequestBody]) AddTags(tags ...string) Route[ResponseBody, RequestBody] {
	r.operation.Tags = append(r.operation.Tags, tags...)
	return r
}

func (r Route[ResponseBody, RequestBody]) RemoveTags(tags ...string) Route[ResponseBody, RequestBody] {
	for _, tag := range tags {
		for i, t := range r.operation.Tags {
			if t == tag {
				r.operation.Tags = slices.Delete(r.operation.Tags, i, i+1)
				break
			}
		}
	}
	return r
}

func (r Route[ResponseBody, RequestBody]) SetDeprecated() Route[ResponseBody, RequestBody] {
	r.operation.Deprecated = true
	return r
}

func UseStd(s *Server, middlewares ...func(http.Handler) http.Handler) {
	Use(s, middlewares...)
}

func Use(s *Server, middlewares ...func(http.Handler) http.Handler) {
	s.middlewares = append(s.middlewares, middlewares...)
}

// Handle registers a standard http handler into the default mux.
// Use this function if you want to use a standard http handler instead of a fuego controller.
func Handle(s *Server, path string, controller http.Handler, middlewares ...func(http.Handler) http.Handler) Route[any, any] {
	return register[any, any](s, http.MethodGet, path, controller, middlewares...)
}

func GetStd(s *Server, path string, controller func(http.ResponseWriter, *http.Request), middlewares ...func(http.Handler) http.Handler) Route[any, any] {
	return RegisterStd(s, http.MethodGet, path, controller, middlewares...)
}

func PostStd(s *Server, path string, controller func(http.ResponseWriter, *http.Request), middlewares ...func(http.Handler) http.Handler) Route[any, any] {
	return RegisterStd(s, http.MethodPost, path, controller, middlewares...)
}

func DeleteStd(s *Server, path string, controller func(http.ResponseWriter, *http.Request), middlewares ...func(http.Handler) http.Handler) Route[any, any] {
	return RegisterStd(s, http.MethodDelete, path, controller, middlewares...)
}

func PutStd(s *Server, path string, controller func(http.ResponseWriter, *http.Request), middlewares ...func(http.Handler) http.Handler) Route[any, any] {
	return RegisterStd(s, http.MethodPut, path, controller, middlewares...)
}

func PatchStd(s *Server, path string, controller func(http.ResponseWriter, *http.Request), middlewares ...func(http.Handler) http.Handler) Route[any, any] {
	return RegisterStd(s, http.MethodPatch, path, controller, middlewares...)
}

// RegisterStd registers a standard http handler into the default mux.
func RegisterStd(s *Server, method string, path string, controller func(http.ResponseWriter, *http.Request), middlewares ...func(http.Handler) http.Handler) Route[any, any] {
	fullPath := method + " " + s.basePath + path

	slog.Debug("registering standard controller " + fullPath)
	route := register[any, any](s, method, path, http.HandlerFunc(controller), middlewares...)

	name, nameWithPath := funcName(controller)
	route.operation.Summary = name
	route.operation.Description = "controller: " + nameWithPath
	route.operation.OperationID = fullPath + ":" + name
	return route
}

func withMiddlewares(controller http.Handler, middlewares ...func(http.Handler) http.Handler) http.Handler {
	for _, middleware := range middlewares {
		controller = middleware(controller)
	}
	return controller
}

// funcName returns the name of a function and the name with package path
func funcName(f interface{}) (name string, nameWithPath string) {
	nameWithPath = strings.TrimSuffix(runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name(), "-fm")
	fullName := strings.Split(nameWithPath, ".")
	return fullName[len(fullName)-1], nameWithPath
}
