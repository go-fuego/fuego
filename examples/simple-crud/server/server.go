package server

import (
	"net/http"

	"simple-crud/controller"
	"simple-crud/static"
	"simple-crud/templates"
	"simple-crud/views"

	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/go-fuego/fuego"
)

type Ressources struct {
	Views views.Ressource
	API   controller.Ressource
}

func (rs Ressources) Setup(
	options ...func(*fuego.Server),
) *fuego.Server {
	serverOptions := []func(*fuego.Server){
		fuego.WithAutoAuth(controller.LoginFunc),
		fuego.WithTemplateFS(templates.FS),
		fuego.WithTemplateGlobs("**/*.html", "**/**/*.html"),
	}

	options = append(serverOptions, options...)

	// Create server with some options
	app := fuego.NewServer(options...)

	rs.API.Security = app.Security

	// Register middlewares (functions that will be executed before AND after the controllers, in the order they are registered)
	// With fuego, you can use any existing middleware that relies on `net/http`, or create your own
	fuego.Use(app, chiMiddleware.Compress(5, "text/html", "text/css", "application/json"))

	fuego.Handle(app, "/static/", http.StripPrefix("/static", static.Handler()))
	fuego.Handle(app, "/manifest.json", static.Handler())
	fuego.Handle(app, "/favicon.ico", static.Handler())

	// Register views (controllers that return HTML pages)
	rs.Views.Routes(fuego.Group(app, "/"))

	// Register API routes (controllers that return JSON)
	rs.API.MountRoutes(fuego.Group(app, "/api"))

	return app
}
