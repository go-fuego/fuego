package main

import (
	"flag"
	"log/slog"
	"os"

	"simple-crud/controller"
	"simple-crud/store"

	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/go-op/op"
	"github.com/lmittmann/tint"
	"github.com/rs/cors"
)

//go:generate sqlc generate

func main() {
	debug := flag.Bool("debug", false, "debug mode")
	flag.Parse()

	logLevel := slog.LevelInfo
	if *debug {
		logLevel = slog.LevelDebug
	}

	slog.SetDefault(slog.New(
		tint.NewHandler(os.Stderr, &tint.Options{
			AddSource:  true,
			Level:      logLevel,
			TimeFormat: "15:04:05",
		}),
	))

	db := store.InitDB("/tmp/recipe.db")

	queries := store.New(db)

	rs := controller.NewRessource(*queries)

	app := op.NewServer()
	op.UseStd(app, cors.Default().Handler)
	op.UseStd(app, chiMiddleware.Compress(5, "text/html", "text/css", "application/json"))

	rs.Routes(app)

	op.Group(app, "/api", func(newS *op.Server) {
		op.Get(newS, "/mounted-route", func(c op.Ctx[any]) (string, error) {
			return "hello", nil
		})

		op.Post(newS, "/mounted-route-post", func(c op.Ctx[any]) (string, error) {
			return "hello", nil
		})

		op.Group(newS, "/mounted-group", func(groupedS *op.Server) {
			op.Get(groupedS, "/mounted-route", func(c op.Ctx[any]) (string, error) {
				return "hello", nil
			})
		})
	})

	app.Run()
}
