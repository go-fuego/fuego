package static

import (
	"embed"
	"net/http"
)

//go:embed *
var StaticFiles embed.FS

// Handler returns a http.Handler that will serve files from
// the given file system.
func Handler() http.Handler {
	return http.FileServer(http.FS(StaticFiles))
}
