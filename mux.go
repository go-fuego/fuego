package op

import (
	"log/slog"
	"net/http"
	"reflect"
	"runtime"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
)

func Group(s *Server, path string, group func(s *Server)) {
	ss := *s
	ss.basePath += path

	group(&ss)
}

type Route[ResponseBody any, RequestBody any] struct {
	operation *openapi3.Operation
}

func Get[T any, B any](s *Server, path string, controller func(Ctx[B]) (T, error)) Route[T, B] {
	return Register[T](s, http.MethodGet, path, controller)
}

func Post[T any, B any](s *Server, path string, controller func(Ctx[B]) (T, error)) Route[T, B] {
	return Register[T](s, http.MethodPost, path, controller)
}

func Delete[T any, B any](s *Server, path string, controller func(Ctx[B]) (T, error)) Route[T, B] {
	return Register[T](s, http.MethodDelete, path, controller)
}

func Put[T any, B any](s *Server, path string, controller func(Ctx[B]) (T, error)) Route[T, B] {
	return Register[T](s, http.MethodPut, path, controller)
}

func Patch[T any, B any](s *Server, path string, controller func(Ctx[B]) (T, error)) Route[T, B] {
	return Register[T](s, http.MethodPatch, path, controller)
}

// Registers route into the default mux.
func Register[T any, B any](s *Server, method string, path string, controller func(Ctx[B]) (T, error)) Route[T, B] {
	fullPath := s.basePath + path
	if isGo1_22 {
		fullPath = method + " " + path
	}
	slog.Debug("registering openapi controller " + fullPath)
	route := register[T, B](s, method, path, httpHandler[T, B](s, controller))

	route.operation.Summary = funcName(controller)
	route.operation.Description = "controller: " + funcPathAndName(controller)
	route.operation.OperationID = fullPath + ":" + funcName(controller)
	return route
}

func register[T any, B any](s *Server, method string, path string, controller func(http.ResponseWriter, *http.Request)) Route[T, B] {
	fullPath := s.basePath + path
	if isGo1_22 {
		fullPath = method + " " + path
	}

	s.mux.Handle(fullPath, withMiddlewares(http.HandlerFunc(controller), s.middlewares...))

	operation, err := RegisterOpenAPIOperation[T, B](s, method, s.basePath+path)
	if err != nil {
		slog.Warn("error documenting openapi operation", "error", err)
	}

	name := funcName(controller)
	operation.Summary = name
	operation.Description = "controller: " + funcPathAndName(controller)
	operation.OperationID = fullPath + ":" + name

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

func (r Route[ResponseBody, RequestBody]) WithTags(tags ...string) Route[ResponseBody, RequestBody] {
	r.operation.Tags = append(r.operation.Tags, tags...)
	return r
}

func (r Route[ResponseBody, RequestBody]) WithDeprecated() Route[ResponseBody, RequestBody] {
	r.operation.Deprecated = true
	return r
}

func (r Route[ResponseBody, RequestBody]) WithQueryParam(name, description string) Route[ResponseBody, RequestBody] {
	parameter := openapi3.NewQueryParameter(name)
	parameter.Description = description
	parameter.Schema = openapi3.NewStringSchema().NewRef()
	r.operation.AddParameter(parameter)
	return r
}

func UseStd(s *Server, middlewares ...func(http.Handler) http.Handler) {
	s.middlewares = append(s.middlewares, middlewares...)
}

func GetStd(s *Server, path string, controller func(http.ResponseWriter, *http.Request)) Route[any, any] {
	return RegisterStd(s, http.MethodGet, path, controller)
}

func PostStd(s *Server, path string, controller func(http.ResponseWriter, *http.Request)) Route[any, any] {
	return RegisterStd(s, http.MethodPost, path, controller)
}

func DeleteStd(s *Server, path string, controller func(http.ResponseWriter, *http.Request)) Route[any, any] {
	return RegisterStd(s, http.MethodDelete, path, controller)
}

func PutStd(s *Server, path string, controller func(http.ResponseWriter, *http.Request)) Route[any, any] {
	return RegisterStd(s, http.MethodPut, path, controller)
}

func PatchStd(s *Server, path string, controller func(http.ResponseWriter, *http.Request)) Route[any, any] {
	return RegisterStd(s, http.MethodPatch, path, controller)
}

// RegisterStd registers a standard http handler into the default mux.
func RegisterStd(s *Server, method string, path string, controller func(http.ResponseWriter, *http.Request)) Route[any, any] {
	fullPath := s.basePath + path
	if isGo1_22 {
		fullPath = method + " " + path
	}
	slog.Debug("registering standard controller " + fullPath)
	return register[any, any](s, method, path, controller)
}

func withMiddlewares(controller http.Handler, middlewares ...func(http.Handler) http.Handler) http.Handler {
	for _, middleware := range middlewares {
		controller = middleware(controller)
	}
	return controller
}

// funcPathAndName returns the path and name of a function.
func funcPathAndName(f interface{}) string {
	return strings.TrimSuffix(runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name(), "-fm")
}

func funcName(f interface{}) string {
	fullName := strings.Split(funcPathAndName(f), ".")
	return fullName[len(fullName)-1]
}
