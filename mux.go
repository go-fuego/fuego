package fuego

import (
	"log/slog"
	"net/http"
	"reflect"
	"runtime"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
)

// Group allows grouping routes under a common path.
// Middlewares are scoped to the group.
// For example:
//
//	s := fuego.NewServer()
//	viewsRoutes := fuego.Group(s, "")
//	apiRoutes := fuego.Group(s, "/api")
//	// Registering a middlewares scoped to /api only
//	fuego.Use(apiRoutes, myMiddleware)
//	// Registering a route under /api/users
//	fuego.Get(apiRoutes, "/users", func(c fuego.ContextNoBody) (ans, error) {
//		return ans{Ans: "users"}, nil
//	})
//	s.Run()
func Group(s *Server, path string, routeOptions ...func(*BaseRoute)) *Server {
	if path == "/" {
		path = ""
	} else if path != "" && path[len(path)-1] == '/' {
		slog.Warn("Group path should not end with a slash.", "path", path+"/", "new", path)
	}

	ss := *s
	newServer := &ss
	newServer.basePath += path

	if autoTag := strings.TrimLeft(path, "/"); !s.disableAutoGroupTags && autoTag != "" {
		newServer.routeOptions = append(s.routeOptions, OptionTags(autoTag))
	}

	newServer.routeOptions = append(newServer.routeOptions, routeOptions...)

	return newServer
}

type Route[ResponseBody any, RequestBody any] struct {
	BaseRoute
}

type BaseRoute struct {
	Operation            *openapi3.Operation // GENERATED OpenAPI operation, do not set manually in Register function. You can change it after the route is registered.
	Method               string              // HTTP method (GET, POST, PUT, PATCH, DELETE)
	Path                 string              // URL path. Will be prefixed by the base path of the server and the group path if any
	Handler              http.Handler        // handler executed for this route
	FullName             string              // namespace and name of the function to execute
	Params               map[string]OpenAPIParam
	Middlewares          []func(http.Handler) http.Handler
	AcceptedContentTypes []string // Content types accepted for the request body. If nil, all content types (*/*) are accepted.
	Hidden               bool     // If true, the route will not be documented in the OpenAPI spec
	DefaultStatusCode    int      // Default status code for the response
	OpenAPI              *OpenAPI // Ref to the whole OpenAPI spec
}

// Capture all methods (GET, POST, PUT, PATCH, DELETE) and register a controller.
func All[ReturnType, Body any, Contexted ctx[Body]](s *Server, path string, controller func(Contexted) (ReturnType, error), options ...func(*BaseRoute)) *Route[ReturnType, Body] {
	return registerFuegoController(s, "", path, controller, options...)
}

func Get[T, B any, Contexted ctx[B]](s *Server, path string, controller func(Contexted) (T, error), options ...func(*BaseRoute)) *Route[T, B] {
	return registerFuegoController(s, http.MethodGet, path, controller, options...)
}

func Post[T, B any, Contexted ctx[B]](s *Server, path string, controller func(Contexted) (T, error), options ...func(*BaseRoute)) *Route[T, B] {
	return registerFuegoController(s, http.MethodPost, path, controller, options...)
}

func Delete[T, B any, Contexted ctx[B]](s *Server, path string, controller func(Contexted) (T, error), options ...func(*BaseRoute)) *Route[T, B] {
	return registerFuegoController(s, http.MethodDelete, path, controller, options...)
}

func Put[T, B any, Contexted ctx[B]](s *Server, path string, controller func(Contexted) (T, error), options ...func(*BaseRoute)) *Route[T, B] {
	return registerFuegoController(s, http.MethodPut, path, controller, options...)
}

func Patch[T, B any, Contexted ctx[B]](s *Server, path string, controller func(Contexted) (T, error), options ...func(*BaseRoute)) *Route[T, B] {
	return registerFuegoController(s, http.MethodPatch, path, controller, options...)
}

// Register registers a controller into the default mux and documents it in the OpenAPI spec.
func Register[T, B any](s *Server, route Route[T, B], controller http.Handler, options ...func(*BaseRoute)) *Route[T, B] {
	for _, o := range options {
		o(&route.BaseRoute)
	}
	route.Handler = controller

	fullPath := s.basePath + route.Path
	if route.Method != "" {
		fullPath = route.Method + " " + fullPath
	}
	slog.Debug("registering controller " + fullPath)

	allMiddlewares := append(s.middlewares, route.Middlewares...)
	s.Mux.Handle(fullPath, withMiddlewares(route.Handler, allMiddlewares...))

	if s.DisableOpenapi || route.Hidden || route.Method == "" {
		return &route
	}

	route.Path = s.basePath + route.Path

	err := route.RegisterOpenAPIOperation(s.OpenAPI)
	if err != nil {
		slog.Warn("error documenting openapi operation", "error", err)
	}

	if route.FullName == "" {
		route.FullName = route.Path
	}

	if route.Operation.Summary == "" {
		route.Operation.Summary = route.NameFromNamespace(camelToHuman)
	}

	route.Operation.Description = "controller: `" + route.FullName + "`\n\n---\n\n" + route.Operation.Description

	if route.Operation.OperationID == "" {
		route.Operation.OperationID = route.Method + "_" + strings.ReplaceAll(strings.ReplaceAll(route.Path, "{", ":"), "}", "")
	}

	return &route
}

func UseStd(s *Server, middlewares ...func(http.Handler) http.Handler) {
	Use(s, middlewares...)
}

func Use(s *Server, middlewares ...func(http.Handler) http.Handler) {
	s.middlewares = append(s.middlewares, middlewares...)
}

// Handle registers a standard HTTP handler into the default mux.
// Use this function if you want to use a standard HTTP handler instead of a Fuego controller.
func Handle(s *Server, path string, controller http.Handler, options ...func(*BaseRoute)) *Route[any, any] {
	return Register(s, Route[any, any]{
		BaseRoute: BaseRoute{
			Path:     path,
			FullName: FuncName(controller),
		},
	}, controller)
}

func AllStd(s *Server, path string, controller func(http.ResponseWriter, *http.Request), options ...func(*BaseRoute)) *Route[any, any] {
	return registerStdController(s, "", path, controller, options...)
}

func GetStd(s *Server, path string, controller func(http.ResponseWriter, *http.Request), options ...func(*BaseRoute)) *Route[any, any] {
	return registerStdController(s, http.MethodGet, path, controller, options...)
}

func PostStd(s *Server, path string, controller func(http.ResponseWriter, *http.Request), options ...func(*BaseRoute)) *Route[any, any] {
	return registerStdController(s, http.MethodPost, path, controller, options...)
}

func DeleteStd(s *Server, path string, controller func(http.ResponseWriter, *http.Request), options ...func(*BaseRoute)) *Route[any, any] {
	return registerStdController(s, http.MethodDelete, path, controller, options...)
}

func PutStd(s *Server, path string, controller func(http.ResponseWriter, *http.Request), options ...func(*BaseRoute)) *Route[any, any] {
	return registerStdController(s, http.MethodPut, path, controller, options...)
}

func PatchStd(s *Server, path string, controller func(http.ResponseWriter, *http.Request), options ...func(*BaseRoute)) *Route[any, any] {
	return registerStdController(s, http.MethodPatch, path, controller, options...)
}

func registerFuegoController[T, B any, Contexted ctx[B]](s *Server, method, path string, controller func(Contexted) (T, error), options ...func(*BaseRoute)) *Route[T, B] {
	route := BaseRoute{
		Method:    method,
		Path:      path,
		Params:    make(map[string]OpenAPIParam),
		FullName:  FuncName(controller),
		Operation: openapi3.NewOperation(),
		OpenAPI:   s.OpenAPI,
	}

	acceptHeaderParameter := openapi3.NewHeaderParameter("Accept")
	acceptHeaderParameter.Schema = openapi3.NewStringSchema().NewRef()
	route.Operation.AddParameter(acceptHeaderParameter)

	for _, o := range append(s.routeOptions, options...) {
		o(&route)
	}

	return Register(s, Route[T, B]{BaseRoute: route}, HTTPHandler(s, controller, &route))
}

func registerStdController(s *Server, method, path string, controller func(http.ResponseWriter, *http.Request), options ...func(*BaseRoute)) *Route[any, any] {
	route := BaseRoute{
		Method:    method,
		Path:      path,
		FullName:  FuncName(controller),
		Operation: openapi3.NewOperation(),
		OpenAPI:   s.OpenAPI,
	}

	for _, o := range append(s.routeOptions, options...) {
		o(&route)
	}

	return Register(s, Route[any, any]{BaseRoute: route}, http.HandlerFunc(controller))
}

func withMiddlewares(controller http.Handler, middlewares ...func(http.Handler) http.Handler) http.Handler {
	for i := len(middlewares) - 1; i >= 0; i-- {
		controller = middlewares[i](controller)
	}
	return controller
}

// FuncName returns the name of a function and the name with package path
func FuncName(f interface{}) string {
	return strings.TrimSuffix(runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name(), "-fm")
}

// NameFromNamespace returns the Route's FullName final string
// delimited by `.`. Essentially getting the name of the function
// and leaving the package path
//
// The output can be further modified with a list of optional
// string manipulation funcs (i.e func(string) string)
func (r Route[T, B]) NameFromNamespace(opts ...func(string) string) string {
	ss := strings.Split(r.FullName, ".")
	name := ss[len(ss)-1]
	for _, o := range opts {
		name = o(name)
	}
	return name
}

// transform camelCase to human-readable string
func camelToHuman(s string) string {
	result := strings.Builder{}
	for i, r := range s {
		if 'A' <= r && r <= 'Z' {
			if i > 0 {
				result.WriteRune(' ')
			}
			result.WriteRune(r + 'a' - 'A') // 'A' -> 'a
		} else {
			result.WriteRune(r)
		}
	}
	return result.String()
}
