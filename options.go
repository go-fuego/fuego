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
	Serialize             func(w http.ResponseWriter, ans any)
	SerializeError        func(w http.ResponseWriter, err error)
	startTime             time.Time
}

// NewServer creates a new server with the given options
func NewServer(options ...func(*Server)) *Server {
	s := &Server{
		mux:  http.NewServeMux(),
		spec: NewOpenAPI(),

		Addr:                  ":8080",
		DisallowUnknownFields: true,
		Serialize:             SendJSON,
		SerializeError:        SendJSONError,
	}

	for _, option := range options {
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

func WithDisallowUnknownFields(b bool) func(*Server) {
	return func(c *Server) { c.DisallowUnknownFields = b }
}

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

func WithoutLogger() func(*Server) {
	return func(c *Server) {
		slog.SetDefault(slog.New(slog.NewTextHandler(io.Discard, nil)))
	}
}
