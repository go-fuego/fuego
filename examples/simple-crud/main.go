package main

import (
	"flag"
	"log/slog"
	"os"

	"simple-crud/controller"
	"simple-crud/store"

	"github.com/go-op/op"
	"github.com/lmittmann/tint"
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
	rs.Routes(app)
	app.Run()
}
