package fuego

import (
	"crypto/tls"
	"errors"
	"fmt"
	"html/template"
	"log/slog"
	"net"
	"net/http"
	"reflect"
	"time"
)

// Run starts the server.
// It is blocking.
// It returns an error if the server could not start (it could not bind to the port for example).
// It also generates the OpenAPI spec and outputs it to a file, the UI, and a handler (if enabled).
func (s *Server) Run() error {
	if err := s.setupDefaultListener(); err != nil {
		return fmt.Errorf("failed to start server: %w", err)
	}
	s.setup()
	return s.Server.Serve(s.listener)
}

// RunTLS starts the server with a TLS listener
// It is blocking.
// It returns an error if the server could not start (it could not bind to the port for example).
// It also generates the OpenAPI spec and outputs it to a file, the UI, and a handler (if enabled).
func (s *Server) RunTLS(certFile, keyFile string) error {
	s.isTLS = true
	if err := s.setupTLSListener(certFile, keyFile); err != nil {
		return fmt.Errorf("failed to start TLS server: %w", err)
	}
	s.setup()
	return s.Server.Serve(s.listener)
}

func (s *Server) setup() {
	go s.OutputOpenAPISpec()
	s.printStartupMessage()

	s.Server.Handler = s.Mux
	if s.corsMiddleware != nil {
		s.Server.Handler = s.corsMiddleware(s.Server.Handler)
	}
}

// setupTLSListener creates a TLS listener if no listener is already configured.
// If a non-TLS listener is already configured, an error is returned.
// Requires valid TLS certificate and key files to establish a secure listener.
// Returns an error if the listener cannot be created or if the provided certificates are invalid.
func (s *Server) setupTLSListener(certFile, keyFile string) error {
	if s.listener != nil && !s.isTLS {
		return errors.New("a non-TLS listener is already configured; cannot set up a TLS listener on the same server")
	}
	if s.listener != nil {
		return errors.New("a TLS listener is already configured; use the Run() method to start the server")
	}
	if certFile == "" || keyFile == "" {
		return errors.New("TLS certificate and key files must be provided to set up a TLS listener")
	}
	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return fmt.Errorf("failed to load TLS certificate and key files (%s, %s): %w", certFile, keyFile, err)
	}
	tlsConfig := &tls.Config{Certificates: []tls.Certificate{cert}}

	listener, err := tls.Listen("tcp", s.Server.Addr, tlsConfig)
	if err != nil {
		return fmt.Errorf("failed to create a TLS listener on address %s: %w", s.Server.Addr, err)
	}
	s.listener = listener
	return nil
}

// setupDefaultListener creates a default (non-TLS) listener if none is already configured.
// If a listener is already set, this method does nothing.
// Returns an error if the listener cannot be created (e.g., address binding issues).
func (s *Server) setupDefaultListener() error {
	if s.listener != nil {
		return nil // Listener already exists, no action needed.
	}
	listener, err := net.Listen("tcp", s.Server.Addr)
	if err != nil {
		return fmt.Errorf("failed to create default listener on %s: %w", s.Server.Addr, err)
	}
	s.listener = listener
	return nil
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
	return s.proto() + "://" + s.Server.Addr
}

// HTTPHandler converts a Fuego controller into a http.HandlerFunc.
// Uses Server for configuration.
// Uses Route for route configuration. Optional.
func HTTPHandler[ReturnType, Body any](s *Server, controller func(c ContextWithBody[Body]) (ReturnType, error), route *BaseRoute) http.HandlerFunc {
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

		ctx := &netHttpContext[Body]{
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
		}

		err := validateParams(*ctx)
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

		if route.DefaultStatusCode != 0 {
			w.WriteHeader(route.DefaultStatusCode)
		}

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
