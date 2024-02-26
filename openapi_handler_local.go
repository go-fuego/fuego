//go:build openapi_ui_local

package fuego

import (
	"net/http"

	httpSwagger "github.com/swaggo/http-swagger"
)

func openApiHandler(specURL string) http.Handler {
	return httpSwagger.Handler(
		httpSwagger.Layout(httpSwagger.BaseLayout),
		httpSwagger.PersistAuthorization(true),
		httpSwagger.URL(specURL), // The url pointing to API definition
	)
}
