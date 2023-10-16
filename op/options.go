package op

type Config struct {
	Addr                  string
	DisallowUnknownFields bool // If true, the server will return an error if the request body contains unknown fields. Useful for quick debugging in development.
}

// Config variable used in all project.
var config = Config{
	Addr:                  ":8080",
	DisallowUnknownFields: true,
}

// Options sets the options for the server.
func Options(options ...func(*Config)) {
	for _, option := range options {
		option(&config)
	}
}

func WithDisallowUnknownFields(b bool) func(*Config) {
	return func(c *Config) { c.DisallowUnknownFields = b }
}

func WithPort(port string) func(*Config) {
	return func(c *Config) { c.Addr = port }
}
