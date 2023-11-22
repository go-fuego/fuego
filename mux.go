package op

import (
	"log/slog"
	"net/http"
	"reflect"
	"runtime"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
)

func Group(s *Server, path string, group func(s *Server)) *Server {
	ss := *s
	newServer := &ss
	newServer.basePath += path

	if group != nil {
		group(newServer)
	}
	return newServer
}

type Route[ResponseBody any, RequestBody any] struct {
	operation *openapi3.Operation
}

func Get[T any, B any](s *Server, path string, controller func(Ctx[B]) (T, error), middlewares ...func(http.Handler) http.Handler) Route[T, B] {
	return Register[T](s, http.MethodGet, path, controller, middlewares...)
}

func Post[T any, B any](s *Server, path string, controller func(Ctx[B]) (T, error), middlewares ...func(http.Handler) http.Handler) Route[T, B] {
	return Register[T](s, http.MethodPost, path, controller, middlewares...)
}

func Delete[T any, B any](s *Server, path string, controller func(Ctx[B]) (T, error), middlewares ...func(http.Handler) http.Handler) Route[T, B] {
	return Register[T](s, http.MethodDelete, path, controller, middlewares...)
}

func Put[T any, B any](s *Server, path string, controller func(Ctx[B]) (T, error), middlewares ...func(http.Handler) http.Handler) Route[T, B] {
	return Register[T](s, http.MethodPut, path, controller, middlewares...)
}

func Patch[T any, B any](s *Server, path string, controller func(Ctx[B]) (T, error), middlewares ...func(http.Handler) http.Handler) Route[T, B] {
	return Register[T](s, http.MethodPatch, path, controller, middlewares...)
}

// Registers route into the default mux.
func Register[T any, B any](s *Server, method string, path string, controller func(Ctx[B]) (T, error), middlewares ...func(http.Handler) http.Handler) Route[T, B] {
	fullPath := s.basePath + path
	if isGo1_22 {
		fullPath = method + " " + path
	}
	slog.Debug("registering openapi controller " + fullPath)
	route := register[T, B](s, method, path, httpHandler[T, B](s, controller), middlewares...)

	name, nameWithPath := funcName(controller)
	route.operation.Summary = name
	route.operation.Description = "controller: " + nameWithPath
	route.operation.OperationID = fullPath + ":" + name
	return route
}

func register[T any, B any](s *Server, method string, path string, controller func(http.ResponseWriter, *http.Request), middlewares ...func(http.Handler) http.Handler) Route[T, B] {
	fullPath := s.basePath + path
	if isGo1_22 {
		fullPath = method + " " + path
	}

	allMiddlewares := append(middlewares, s.middlewares...)
	s.mux.Handle(fullPath, withMiddlewares(http.HandlerFunc(controller), allMiddlewares...))

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

func (r Route[ResponseBody, RequestBody]) AddTags(tags ...string) Route[ResponseBody, RequestBody] {
	r.operation.Tags = append(r.operation.Tags, tags...)
	return r
}

func (r Route[ResponseBody, RequestBody]) SetDeprecated() Route[ResponseBody, RequestBody] {
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
	fullPath := s.basePath + path
	if isGo1_22 {
		fullPath = method + " " + path
	}
	slog.Debug("registering standard controller " + fullPath)
	route := register[any, any](s, method, path, controller, middlewares...)

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
