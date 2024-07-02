package fuego

import (
	"errors"
	"fmt"
	"html/template"
	"log/slog"
	"net/http"
	"reflect"
	"strings"
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

// initializes any Context type with the base ContextNoBody context.
//
//	var ctx ContextWithBody[any] // does not work because it will create a ContextWithBody[any] with a nil value
func initContext[Contextable ctx[Body], Body any](baseContext ContextNoBody) Contextable {
	var c Contextable

	switch any(c).(type) {
	case ContextNoBody:
		return any(baseContext).(Contextable)
	case *ContextNoBody:
		return any(&baseContext).(Contextable)
	case *ContextWithBody[Body]:
		return any(&ContextWithBody[Body]{
			ContextNoBody: baseContext,
		}).(Contextable)
	default:
		panic("unknown type")
	}
}

// HTTPHandler converts a Fuego controller into a http.HandlerFunc.
func HTTPHandler[ReturnType, Body any, Contextable ctx[Body]](s *Server, controller func(c Contextable) (ReturnType, error)) http.HandlerFunc {
	baseContext := *new(Contextable)
	if reflect.TypeOf(baseContext) == nil {
		slog.Info(fmt.Sprintf("context is nil: %v %T", baseContext, baseContext))
		panic("ctx must be provided as concrete type (not interface). ContextNoBody, ContextWithBody[any], ContextFull[any, any], ContextWithQueryParams[any] are supported")
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
		ctx := initContext[Contextable](ContextNoBody{
			Req: r,
			Res: w,
			readOptions: readOptions{
				DisallowUnknownFields: s.DisallowUnknownFields,
				MaxBodySize:           s.maxBodySize,
			},
			fs:        s.fs,
			templates: templates,
		})

		timeController := time.Now()
		w.Header().Set("Server-Timing", Timing{"fuegoReqInit", timeController.Sub(timeCtxInit), ""}.String())

		// CONTROLLER
		ans, err := controller(ctx)
		if err != nil {
			err = s.ErrorHandler(err)
			SendError(w, r, err)
			return
		}
		timeAfterController := time.Now()
		w.Header().Add("Server-Timing", Timing{"controller", timeAfterController.Sub(timeController), ""}.String())

		if reflect.TypeOf(ans) == nil {
			return
		}

		// TRANSFORM OUT
		timeTransformOut := time.Now()
		ans, err = transformOut(r.Context(), ans)
		if err != nil {
			err = s.ErrorHandler(err)
			SendError(w, r, err)
			return
		}
		timeAfterTransformOut := time.Now()
		w.Header().Add("Server-Timing", Timing{"transformOut", timeAfterTransformOut.Sub(timeTransformOut), "transformOut"}.String())

		// SERIALIZATION
		switch parseAcceptHeader(r.Header.Get("Accept"), ans) {
		case "application/xml":
			SendXML(w, ans)
		case "text/html":
			err = SendHTML(r.Context(), w, ans)
		case "text/plain":
			err = SendText(w, ans)
		case "application/json":
			SendJSON(w, ans)
		default:
			SendError(w, r, errors.New("unsupported Accept header"))
		}
		w.Header().Add("Server-Timing", Timing{"serialize", time.Since(timeAfterTransformOut), ""}.String())

		// FINAL ERROR HANDLING
		if err != nil {
			err = s.ErrorHandler(err)
			SendError(w, r, err)
		}
	}
}

func parseAcceptHeader(accept string, ans any) string {
	if strings.Index(accept, ",") > 0 {
		accept = accept[:strings.Index(accept, ",")]
	}
	if accept == "*/*" {
		accept = ""
	}
	if accept == "" {
		accept = InferAcceptHeaderFromType(ans)
	}
	return accept
}
