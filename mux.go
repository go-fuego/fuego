package fuego

import (
	"log/slog"
	"net/http"
	"reflect"
	"runtime"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
)

// Group allows to group routes under a common path.
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
func Group(s *Server, path string) *Server {
	if path == "/" {
		path = ""
	} else if path != "" && path[len(path)-1] == '/' {
		slog.Warn("Group path should not end with a slash.", "path", path+"/", "new", path)
	}

	ss := *s
	newServer := &ss
	newServer.basePath += path
	newServer.groupTag = strings.TrimLeft(path, "/")

	return newServer
}

type Route[ResponseBody any, RequestBody any] struct {
	Operation *openapi3.Operation
	Method    string       // HTTP method (GET, POST, PUT, PATCH, DELETE)
	Path      string       // URL path. Will be prefixed by the base path of the server and the group path if any
	Handler   http.Handler // handler executed for this route
	FullName  string       // namespace and name of the function to execute
}

// Capture all methods (GET, POST, PUT, PATCH, DELETE) and register a controller.
func All[T, B any, Contexted ctx[B]](s *Server, path string, controller func(Contexted) (T, error), middlewares ...func(http.Handler) http.Handler) Route[T, B] {
	return Register(s, Route[T, B]{
		Path:     path,
		FullName: FuncName(controller),
	}, HTTPHandler(s, controller), middlewares...)
}

func Get[T, B any, Contexted ctx[B]](s *Server, path string, controller func(Contexted) (T, error), middlewares ...func(http.Handler) http.Handler) Route[T, B] {
	return Register(s, Route[T, B]{
		Method:   http.MethodGet,
		Path:     path,
		FullName: FuncName(controller),
	}, HTTPHandler(s, controller), middlewares...)
}

func Post[T, B any, Contexted ctx[B]](s *Server, path string, controller func(Contexted) (T, error), middlewares ...func(http.Handler) http.Handler) Route[T, B] {
	return Register(s, Route[T, B]{
		Method:   http.MethodPost,
		Path:     path,
		FullName: FuncName(controller),
	}, HTTPHandler(s, controller), middlewares...)
}

func Delete[T, B any, Contexted ctx[B]](s *Server, path string, controller func(Contexted) (T, error), middlewares ...func(http.Handler) http.Handler) Route[T, B] {
	return Register(s, Route[T, B]{
		Method: http.MethodDelete,
		Path:   path,

		FullName: FuncName(controller),
	}, HTTPHandler(s, controller), middlewares...)
}

func Put[T, B any, Contexted ctx[B]](s *Server, path string, controller func(Contexted) (T, error), middlewares ...func(http.Handler) http.Handler) Route[T, B] {
	return Register(s, Route[T, B]{
		Method:   http.MethodPut,
		Path:     path,
		FullName: FuncName(controller),
	}, HTTPHandler(s, controller), middlewares...)
}

func Patch[T, B any, Contexted ctx[B]](s *Server, path string, controller func(Contexted) (T, error), middlewares ...func(http.Handler) http.Handler) Route[T, B] {
	return Register(s, Route[T, B]{
		Method:   http.MethodPatch,
		Path:     path,
		FullName: FuncName(controller),
	}, HTTPHandler(s, controller), middlewares...)
}

// Register registers a controller into the default mux and documents it in the OpenAPI spec.
func Register[T, B any](s *Server, route Route[T, B], controller http.Handler, middlewares ...func(http.Handler) http.Handler) Route[T, B] {
	route.Handler = controller

	fullPath := s.basePath + route.Path
	if route.Method != "" {
		fullPath = route.Method + " " + fullPath
	}
	slog.Debug("registering controller " + fullPath)

	allMiddlewares := append(s.middlewares, middlewares...)
	s.Mux.Handle(fullPath, withMiddlewares(route.Handler, allMiddlewares...))

	if s.DisableOpenapi || route.Method == "" {
		return route
	}

	var err error
	route.Operation, err = RegisterOpenAPIOperation[T, B](s, route.Method, s.basePath+route.Path)
	if err != nil {
		slog.Warn("error documenting openapi operation", "error", err)
	}

	if route.FullName == "" {
		route.FullName = route.Path
	}

	route.Operation.Summary = route.NameFromNamespace(camelToHuman)
	route.Operation.Description = "controller: `" + route.FullName + "`\n\n---\n\n"
	route.Operation.OperationID = route.Method + " " + s.basePath + route.Path + ":" + route.NameFromNamespace()

	return route
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
	return Register(s, Route[any, any]{
		Path:     path,
		FullName: FuncName(controller),
	}, controller, middlewares...)
}

func AllStd(s *Server, path string, controller func(http.ResponseWriter, *http.Request), middlewares ...func(http.Handler) http.Handler) Route[any, any] {
	return Register(s, Route[any, any]{
		Path:     path,
		FullName: FuncName(controller),
	}, http.HandlerFunc(controller), middlewares...)
}

func GetStd(s *Server, path string, controller func(http.ResponseWriter, *http.Request), middlewares ...func(http.Handler) http.Handler) Route[any, any] {
	return Register(s, Route[any, any]{
		Method:   http.MethodGet,
		Path:     path,
		FullName: FuncName(controller),
	}, http.HandlerFunc(controller), middlewares...)
}

func PostStd(s *Server, path string, controller func(http.ResponseWriter, *http.Request), middlewares ...func(http.Handler) http.Handler) Route[any, any] {
	return Register(s, Route[any, any]{
		Method:   http.MethodPost,
		Path:     path,
		FullName: FuncName(controller),
	}, http.HandlerFunc(controller), middlewares...)
}

func DeleteStd(s *Server, path string, controller func(http.ResponseWriter, *http.Request), middlewares ...func(http.Handler) http.Handler) Route[any, any] {
	return Register(s, Route[any, any]{
		Method:   http.MethodDelete,
		Path:     path,
		FullName: FuncName(controller),
	}, http.HandlerFunc(controller), middlewares...)
}

func PutStd(s *Server, path string, controller func(http.ResponseWriter, *http.Request), middlewares ...func(http.Handler) http.Handler) Route[any, any] {
	return Register(s, Route[any, any]{
		Method:   http.MethodPut,
		Path:     path,
		FullName: FuncName(controller),
	}, http.HandlerFunc(controller), middlewares...)
}

func PatchStd(s *Server, path string, controller func(http.ResponseWriter, *http.Request), middlewares ...func(http.Handler) http.Handler) Route[any, any] {
	return Register(s, Route[any, any]{
		Method:   http.MethodPatch,
		Path:     path,
		FullName: FuncName(controller),
	}, http.HandlerFunc(controller), middlewares...)
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

// transform camelCase to human readable string
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
