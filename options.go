package op

import (
	"io"
	"log/slog"
	"net/http"
	"runtime"
	"strings"
	"time"

	"github.com/getkin/kin-openapi/openapi3"
)

var isGo1_22 = strings.TrimPrefix(runtime.Version(), "devel ") >= "go1.22"

type Server struct {
	middlewares []func(http.Handler) http.Handler
	mux         *http.ServeMux
	basePath    string

	spec openapi3.T

	Addr                  string
	DisallowUnknownFields bool // If true, the server will return an error if the request body contains unknown fields. Useful for quick debugging in development.
	maxBodySize           int64
	Serialize             func(w http.ResponseWriter, ans any)   // Used to serialize the response. Defaults to [SendJSON].
	SerializeError        func(w http.ResponseWriter, err error) // Used to serialize the error response. Defaults to [SendJSONError].
	ErrorHandler          func(err error) error                  // Used to transform any error into a unified error type structure with status code. Defaults to [ErrorHandler]
	startTime             time.Time
}

// NewServer creates a new server with the given options.
// For example:
//
//	app := op.NewServer(
//		op.WithPort(":8080"),
//		op.WithoutLogger(),
//	)
//
// Option all begin with `With`.
// Some default options are set in the function body.
func NewServer(options ...func(*Server)) *Server {
	s := &Server{
		mux:  http.NewServeMux(),
		spec: NewOpenAPI(),
	}

	defaultOptions := [...]func(*Server){
		WithPort(":8080"),
		WithDisallowUnknownFields(true),
		WithSerializer(SendJSON),
		WithErrorSerializer(SendJSONError),
		WithErrorHandler(ErrorHandler),
	}

	for _, option := range append(defaultOptions[:], options...) {
		option(s)
	}

	if !isGo1_22 {
		slog.Warn(
			"Please upgrade to Go >= 1.22. " +
				"You are running " + runtime.Version() + ": " +
				"you cannot use path params nor register routes with the same path but different methods. ")
	}

	s.startTime = time.Now()

	return s
}

// WithDisallowUnknownFields sets the DisallowUnknownFields option.
// If true, the server will return an error if the request body contains unknown fields.
// Useful for quick debugging in development.
// Defaults to true.
func WithDisallowUnknownFields(b bool) func(*Server) {
	return func(c *Server) { c.DisallowUnknownFields = b }
}

// WithPort sets the port of the server. For example, ":8080".
func WithPort(port string) func(*Server) {
	return func(c *Server) { c.Addr = port }
}

func WithXML() func(*Server) {
	return func(c *Server) {
		c.Serialize = SendXML
		c.SerializeError = SendXMLError
	}
}

func WithHandler(handler slog.Handler) func(*Server) {
	return func(c *Server) {
		if handler != nil {
			slog.SetDefault(slog.New(handler))
		}
	}
}

func WithSerializer(serializer func(w http.ResponseWriter, ans any)) func(*Server) {
	return func(c *Server) { c.Serialize = serializer }
}

func WithErrorSerializer(serializer func(w http.ResponseWriter, err error)) func(*Server) {
	return func(c *Server) { c.SerializeError = serializer }
}

func WithErrorHandler(errorHandler func(err error) error) func(*Server) {
	return func(c *Server) { c.ErrorHandler = errorHandler }
}

// WithoutLogger disables the default logger.
func WithoutLogger() func(*Server) {
	return func(c *Server) {
		slog.SetDefault(slog.New(slog.NewTextHandler(io.Discard, nil)))
	}
}
