package basicauth

import (
	"net/http"

	"github.com/go-fuego/fuego"
)

type Config struct {
	Username string
	Password string
	AllowGet bool // Allow GET requests without auth
}

// Basic auth middleware
func New(config Config) func(http.Handler) http.Handler {
	if config.Username == "" {
		panic("basicauth: username is required")
	}
	if config.Password == "" {
		panic("basicauth: password is required")
	}

	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodGet && config.AllowGet {
				h.ServeHTTP(w, r)
				return
			}

			user, pass, ok := r.BasicAuth()

			if ok && user == config.Username && pass == config.Password {
				h.ServeHTTP(w, r)
				return
			}

			err := fuego.HTTPError{
				Title:  "unauthorized access",
				Detail: "wrong username or password",
				Status: http.StatusUnauthorized,
			}

			w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
			fuego.SendJSONError(w, err)
		})
	}
}
