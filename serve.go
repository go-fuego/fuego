package fuego

import (
	"context"
	"html/template"
	"log/slog"
	"net"
	"net/http"
	"reflect"
	"time"

	"github.com/getkin/kin-openapi/openapi3"
)

// Run starts the server.
// It is blocking.
// It returns an error if the server could not start (it could not bind to the port for example).
// It also generates the OpenAPI spec and outputs it to a file, the UI, and a handler (if enabled).
func (s *Server) Run() error {
	return s.RunContext(context.Background())
}

// RunContext runs [Run] but with Context.
// When context is canceled the server is shutdown.
func (s *Server) RunContext(ctx context.Context) error {
	if err := s.setup(); err != nil {
		return err
	}
	return s.serveWithContext(ctx, func() error {
		return s.Serve(s.listener)
	})
}

// RunTLS starts the server with a TLS listener.
// It is blocking.
// It returns an error if the server could not start (it could not bind to the port for example).
// It also generates the OpenAPI spec and outputs it to a file, the UI, and a handler (if enabled).
func (s *Server) RunTLS(certFile, keyFile string) error {
	return s.RunTLSContext(context.Background(), certFile, keyFile)
}

// RunTLSContext runs [RunTLS] but with Context.
// When context is canceled the server is shutdown.
func (s *Server) RunTLSContext(ctx context.Context, certFile, keyFile string) error {
	s.isTLS = true
	if err := s.setup(); err != nil {
		return err
	}
	return s.serveWithContext(ctx, func() error {
		return s.ServeTLS(s.listener, certFile, keyFile)
	})
}

// serveWithContext runs the server and shuts it down when the context is canceled.
func (s *Server) serveWithContext(ctx context.Context, serve func() error) error {
	errCh := make(chan error, 1)
	go func() { errCh <- serve() }()

	select {
	case err := <-errCh:
		return err
	case <-ctx.Done():
		return s.Shutdown(context.WithoutCancel(ctx))
	}
}

func (s *Server) setup() error {
	if err := s.setupDefaultListener(); err != nil {
		return err
	}
	if !s.OpenAPI.Config.DisableDefaultServer {
		s.OpenAPI.Description().Servers = append(s.OpenAPI.Description().Servers, &openapi3.Server{
			URL:         s.url(),
			Description: "local server",
		})
	}
	go s.OutputOpenAPISpec()
	s.RegisterOpenAPIRoutes(s)
	s.printStartupMessage()

	s.Handler = s.Mux

	for _, middleware := range s.globalMiddlewares {
		s.Handler = middleware(s.Handler)
	}

	return nil
}

func (s *Server) setupDefaultListener() error {
	if s.listener != nil {
		s.Addr = s.listener.Addr().String()
		return nil
	}
	listener, err := net.Listen("tcp", s.Addr)
	s.listener = listener
	return err
}

func (s *Server) printStartupMessage() {
	if !s.disableStartupMessages {
		elapsed := time.Since(s.startTime)
		slog.Debug("Server started in "+elapsed.String(), "info", "time between since server creation (fuego.NewServer) and server startup (fuego.Run). Depending on your implementation, there might be things that do not depend on fuego slowing start time")
		slog.Info("Server running ✅ on "+s.url(), "started in", elapsed.String())
	}
}

func (s *Server) proto() string {
	if s.isTLS {
		return "https"
	}
	return "http"
}

func (s *Server) url() string {
	return s.proto() + "://" + s.Addr
}

// HTTPHandler converts a Fuego controller into a http.HandlerFunc.
// Uses Server for configuration.
// Uses Route for route configuration. Optional.
func HTTPHandler[ReturnType, Body, Params any](s *Server, controller func(c Context[Body, Params]) (ReturnType, error), route BaseRoute) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var templates *template.Template
		if s.template != nil {
			templates = template.Must(s.template.Clone())
		}

		// CONTEXT INITIALIZATION
		ctx := NewNetHTTPContext[Body, Params](route, w, r, readOptions{
			DisallowUnknownFields: s.DisallowUnknownFields,
			MaxBodySize:           s.maxBodySize,
		})
		ctx.serializer = s.Serialize
		ctx.errorSerializer = s.SerializeError
		ctx.fs = s.fs
		ctx.templates = templates

		Flow(s.Engine, ctx, controller)
	}
}

// ContextFlowable contains the logic for the flow of a Fuego controller.
// Extends [ContextWithBody] with methods not exposed in the Controllers.
type ContextFlowable[B, P any] interface {
	Context[B, P]

	// SetDefaultStatusCode sets the status code of the response defined in the options.
	SetDefaultStatusCode()
	// Serialize serializes the given data to the response.
	Serialize(data any) error
	// SerializeError serializes the given error to the response.
	SerializeError(err error)
}

// Flow is generic handler for Fuego controllers.
func Flow[B, T, P any](s *Engine, ctx ContextFlowable[B, P], controller func(c Context[B, P]) (T, error)) {
	ctx.SetHeader("X-Powered-By", "Fuego")

	timeCtxInit := time.Now()

	// PARAMS VALIDATION
	err := ValidateParams(ctx)
	if err != nil {
		ctx.SetHeader("Trailer", "Server-Timing")
		err = s.ErrorHandler(ctx, err)
		ctx.SerializeError(err)
		return
	}

	timeController := time.Now()
	ctx.SetHeader("Server-Timing", Timing{"fuegoReqInit", "", timeController.Sub(timeCtxInit)}.String())

	// CONTROLLER
	ans, err := controller(ctx)

	if !isNilError(err) {
		ctx.SetHeader("Trailer", "Server-Timing")
		err = s.ErrorHandler(ctx, err)
		ctx.SerializeError(err)
		return
	}
	ctx.SetHeader("Server-Timing", Timing{"controller", "", time.Since(timeController)}.String())

	ctx.SetDefaultStatusCode()

	if reflect.TypeOf(ans) == nil {
		return
	}

	ctx.SetHeader("Trailer", "Server-Timing")

	// TRANSFORM OUT
	timeTransformOut := time.Now()
	ans, err = transformOut(ctx.Context(), ans)
	if err != nil {
		err = s.ErrorHandler(ctx, err)
		ctx.SerializeError(err)
		return
	}
	timeAfterTransformOut := time.Now()
	ctx.SetHeader("Server-Timing", Timing{"transformOut", "transformOut", timeAfterTransformOut.Sub(timeTransformOut)}.String())

	// SERIALIZATION
	err = ctx.Serialize(ans)
	if err != nil {
		err = s.ErrorHandler(ctx, err)
		ctx.SerializeError(err)
	}
	ctx.SetHeader("Server-Timing", Timing{"serialize", "", time.Since(timeAfterTransformOut)}.String())
}

// check if err isNil. If error is of kind pointer
// check if that is nil.
func isNilError(err error) bool {
	if err == nil {
		return true
	}
	v := reflect.ValueOf(err)
	if v.Kind() == reflect.Ptr {
		return v.IsNil()
	}
	return false
}
