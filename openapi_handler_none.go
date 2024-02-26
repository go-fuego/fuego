//go:build openapi_ui_none

package fuego

import (
	"fmt"
	"net/http"
)

func openApiHandler(specURL string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		_, _ = w.Write([]byte(fmt.Sprintf("OpenAPI specification: %s\nNo UI available because openapi_ui_none build tag was set.", specURL)))
	})
}
