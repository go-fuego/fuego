package fuego

import (
	"fmt"
	"html/template"
	"log/slog"
	"net/http"
	"reflect"
	"time"
)

// Run starts the server.
// It is blocking.
// It returns an error if the server could not start (it could not bind to the port for example).
// It also generates the OpenAPI spec and outputs it to a file, the UI, and a handler (if enabled).
func (s *Server) Run() error {
	s.setup()
	return s.Server.ListenAndServe()
}

// RunTLS starts the server with a TLS listener
// It is blocking.
// It returns an error if the server could not start (it could not bind to the port for example).
// It also generates the OpenAPI spec and outputs it to a file, the UI, and a handler (if enabled).
func (s *Server) RunTLS(certFile, keyFile string) error {
	s.isTLS = true

	s.setup()
	return s.Server.ListenAndServeTLS(certFile, keyFile)
}

func (s *Server) setup() {
	go s.OutputOpenAPISpec()
	s.printStartupMessage()

	s.Server.Handler = s.Mux
	if s.corsMiddleware != nil {
		s.Server.Handler = s.corsMiddleware(s.Server.Handler)
	}
}

func (s *Server) printStartupMessage() {
	if !s.disableStartupMessages {
		elapsed := time.Since(s.startTime)
		slog.Debug("Server started in "+elapsed.String(), "info", "time between since server creation (fuego.NewServer) and server startup (fuego.Run). Depending on your implementation, there might be things that do not depend on fuego slowing start time")
		slog.Info("Server running ✅ on "+s.proto()+"://"+s.Server.Addr, "started in", elapsed.String())
	}
}

func (s *Server) proto() string {
	if s.isTLS {
		return "https"
	}
	return "http"
}

func validateParams(c ContextNoBody) error {
	for k, param := range c.params {
		if param.Default != nil {
			// skip: param has a default
			continue
		}

		if param.Required {
			switch param.Type {
			case QueryParamType:
				if !c.urlValues.Has(k) {
					err := fmt.Errorf("%s is a required query param", k)
					return BadRequestError{
						Title:  "Query Param Not Found",
						Err:    err,
						Detail: "cannot parse request parameter: " + err.Error(),
					}
				}
			case HeaderParamType:
				if !c.HasHeader(k) {
					err := fmt.Errorf("%s is a required header", k)
					return BadRequestError{
						Title:  "Header Not Found",
						Err:    err,
						Detail: "cannot parse request parameter: " + err.Error(),
					}
				}
			}
		}
	}

	return nil
}

// initializes any Context type with the base ContextNoBody context.
//
//	var ctx ContextWithBody[any] // does not work because it will create a ContextWithBody[any] with a nil value
func initContext[Contextable ctx[Body], Body any](baseContext ContextNoBody) (Contextable, error) {
	var c Contextable

	err := validateParams(baseContext)
	if err != nil {
		return c, err
	}

	switch any(c).(type) {
	case ContextNoBody:
		return any(baseContext).(Contextable), nil
	case *ContextNoBody:
		return any(&baseContext).(Contextable), nil
	case *ContextWithBody[Body]:
		return any(&ContextWithBody[Body]{
			ContextNoBody: baseContext,
		}).(Contextable), nil
	default:
		panic("unknown type")
	}
}

// HTTPHandler converts a Fuego controller into a http.HandlerFunc.
// Uses Server for configuration.
// Uses Route for route configuration. Optional.
func HTTPHandler[ReturnType, Body any, Contextable ctx[Body]](s *Server, controller func(c Contextable) (ReturnType, error), route *BaseRoute) http.HandlerFunc {
	// Just a check, not used at request time
	baseContext := *new(Contextable)
	if reflect.TypeOf(baseContext) == nil {
		slog.Info(fmt.Sprintf("context is nil: %v %T", baseContext, baseContext))
		panic("ctx must be provided as concrete type (not interface). ContextNoBody, ContextWithBody[any], ContextFull[any, any], ContextWithQueryParams[any] are supported")
	}

	if route == nil {
		route = &BaseRoute{}
	}

	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Powered-By", "Fuego")
		w.Header().Set("Trailer", "Server-Timing")

		// CONTEXT INITIALIZATION
		timeCtxInit := time.Now()
		var templates *template.Template
		if s.template != nil {
			templates = template.Must(s.template.Clone())
		}

		ctx, err := initContext[Contextable](ContextNoBody{
			Req: r,
			Res: w,
			readOptions: readOptions{
				DisallowUnknownFields: s.DisallowUnknownFields,
				MaxBodySize:           s.maxBodySize,
			},
			fs:        s.fs,
			templates: templates,
			params:    route.Params,
			urlValues: r.URL.Query(),
		})
		if err != nil {
			err = s.ErrorHandler(err)
			s.SerializeError(w, r, err)
			return
		}

		timeController := time.Now()
		w.Header().Set("Server-Timing", Timing{"fuegoReqInit", timeController.Sub(timeCtxInit), ""}.String())

		// CONTROLLER
		ans, err := controller(ctx)
		if err != nil {
			err = s.ErrorHandler(err)
			s.SerializeError(w, r, err)
			return
		}
		w.Header().Add("Server-Timing", Timing{"controller", time.Since(timeController), ""}.String())

		if reflect.TypeOf(ans) == nil {
			return
		}

		// TRANSFORM OUT
		timeTransformOut := time.Now()
		ans, err = transformOut(r.Context(), ans)
		if err != nil {
			err = s.ErrorHandler(err)
			s.SerializeError(w, r, err)
			return
		}
		timeAfterTransformOut := time.Now()
		w.Header().Add("Server-Timing", Timing{"transformOut", timeAfterTransformOut.Sub(timeTransformOut), "transformOut"}.String())

		// SERIALIZATION
		err = s.Serialize(w, r, ans)
		if err != nil {
			err = s.ErrorHandler(err)
			s.SerializeError(w, r, err)
		}
		w.Header().Add("Server-Timing", Timing{"serialize", time.Since(timeAfterTransformOut), ""}.String())
	}
}
