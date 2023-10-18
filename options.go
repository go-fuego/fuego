package op

import (
	"net/http"
)

type Server struct {
	middlewares    []func(http.Handler) http.Handler
	mux            *http.ServeMux
	Config         Config
	Serialize      func(w http.ResponseWriter, ans any)
	SerializeError func(w http.ResponseWriter, err error)
}

type Config struct {
	Addr                  string
	DisallowUnknownFields bool // If true, the server will return an error if the request body contains unknown fields. Useful for quick debugging in development.
}

// NewServer creates a new server with the given options
func NewServer(options ...func(*Config)) *Server {
	s := &Server{
		mux:            http.NewServeMux(),
		Serialize:      SendJSON,
		SerializeError: SendJSONError,
		Config: Config{
			Addr:                  ":8080",
			DisallowUnknownFields: true,
		},
	}

	for _, option := range options {
		option(&s.Config)
	}

	return s
}

func WithDisallowUnknownFields(b bool) func(*Config) {
	return func(c *Config) { c.DisallowUnknownFields = b }
}

func WithPort(port string) func(*Config) {
	return func(c *Config) { c.Addr = port }
}

func WithXML() func(*Server) {
	return func(s *Server) {
		s.Serialize = SendXML
		s.SerializeError = SendXMLError
	}
}
