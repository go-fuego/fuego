package server

import (
	"net/http"

	"github.com/rs/cors"

	"github.com/go-fuego/fuego"
	"github.com/go-fuego/fuego/examples/full-app-gourmet/handler"
	"github.com/go-fuego/fuego/examples/full-app-gourmet/otel"
	"github.com/go-fuego/fuego/examples/full-app-gourmet/static"
	"github.com/go-fuego/fuego/examples/full-app-gourmet/templates"
	"github.com/go-fuego/fuego/option"
)

type Resources struct {
	HandlersResources handler.Resource
}

func cache(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "public, max-age=600")

		h.ServeHTTP(w, r)
	})
}

func (rs Resources) Setup(
	options ...fuego.ServerOption,
) *fuego.Server {
	serverOptions := []fuego.ServerOption{
		fuego.WithTemplateFS(templates.FS),
		fuego.WithTemplateGlobs("**/*.html", "**/**/*.html"),
		fuego.WithGlobalMiddlewares(cors.New(cors.Options{
			AllowOriginFunc:  func(origin string) bool { return true },
			AllowedHeaders:   []string{"*"},
			AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
			AllowCredentials: true,
			MaxAge:           300,
		}).Handler),
		fuego.WithRouteOptions(
			option.AddResponse(http.StatusForbidden, "Forbidden", fuego.Response{Type: fuego.HTTPError{}}),
			option.DefaultResponse("Forbidden", fuego.Response{Type: fuego.HTTPError{}}),
		),
		fuego.WithEngineOptions(
			fuego.WithErrorHandler(customErrorHandler),
			fuego.WithOpenAPIConfig(fuego.OpenAPIConfig{
				PrettyFormatJSON: true,
			}),
			fuego.WithMiddlewareConfig(fuego.MiddlewareConfig{
				ShortMiddlewaresPaths: true,
			}),
		),
	}

	options = append(serverOptions, options...)

	// Create server with some options
	app := fuego.NewServer(options...)

	app.OpenAPI.Description().Info.Title = "Gourmet API"

	// Register middlewares (functions that will be executed before AND after the controllers, in the order they are registered)
	// With fuego, you can use any existing middleware that relies on `net/http`, or create your own

	// Add OpenTelemetry observability middleware first to capture all requests with metrics and traces
	fuego.Use(app, otel.HTTPObservabilityMiddleware)

	fuego.Handle(app, "/static/", http.StripPrefix("/static", static.Handler()), option.Middleware(cache))

	fuego.Use(app,
		TokenToContext(rs.HandlersResources.Security, TokenFromCookie, fuego.TokenFromHeader),
	)

	// Register views (controllers that return HTML pages)
	rs.HandlersResources.Routes(app)

	return app
}
