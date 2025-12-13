package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/joho/godotenv"
	"github.com/lithammer/shortuuid/v4"
	"github.com/lmittmann/tint"
	"golang.org/x/term"

	"github.com/go-fuego/fuego"
	"github.com/go-fuego/fuego/examples/full-app-gourmet/handler"
	"github.com/go-fuego/fuego/examples/full-app-gourmet/otel"
	"github.com/go-fuego/fuego/examples/full-app-gourmet/server"
	"github.com/go-fuego/fuego/examples/full-app-gourmet/store"
)

//go:generate sqlc generate

func main() {
	// Load .env.local then .env files
	err := godotenv.Load(".env.local", ".env")
	if err != nil {
		wd, _ := os.Getwd()
		slog.Warn(fmt.Sprintf("Error loading .env files: %s in dir %s", err, wd))
	}

	// Flags
	port := flag.Int("port", 8083, "port to listen to")
	dbPath := flag.String("db", "data/recipe.db", "path to database file")
	debug := flag.Bool("debug", false, "debug mode")
	jsonLogs := flag.Bool("json-logs", !term.IsTerminal(int(os.Stderr.Fd())), "use JSON logging format (auto-enabled in non-TTY environments)")
	flag.Parse()

	logLevel := slog.LevelInfo
	if *debug {
		logLevel = slog.LevelDebug
	}

	// Set logger based on environment
	var logHandler slog.Handler
	if *jsonLogs {
		logHandler = slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{
			AddSource: true,
			Level:     logLevel,
		})
	} else {
		logHandler = tint.NewHandler(os.Stderr, &tint.Options{
			AddSource:  true,
			Level:      logLevel,
			TimeFormat: "15:04:05",
		})
	}
	slog.SetDefault(slog.New(logHandler))

	// Connect to database
	db := store.InitDB(*dbPath)

	store := store.New(db)
	security := fuego.NewSecurity()
	security.ExpiresInterval = handler.LoginExpirationTime

	// Create resources that will be available in handlers
	viewsResources := handler.Resource{
		RecipesQueries:     store,
		IngredientsQueries: store,
		DosingQueries:      store,
		UsersQueries:       store,
		FavoritesQueries:   store,
		Security:           security,
		HotReload:          *debug,
	}

	rs := server.Resources{
		HandlersResources: viewsResources,
	}

	app := rs.Setup(
		fuego.WithAddr(fmt.Sprintf("%s:%d", os.Getenv("HOST"), *port)),
		fuego.WithLoggingMiddleware(fuego.LoggingConfig{
			RequestIDFunc: func() string { return shortuuid.New() },
		}),
	)

	app.OpenAPI.Description().Servers = append(app.OpenAPI.Description().Servers, &openapi3.Server{
		URL:         os.Getenv("PUBLIC_URL"),
		Description: "Production server",
	})

	// Initialize OpenTelemetry
	ctx := context.Background()

	// Initialize metrics
	shutdownMetrics, err := otel.InitMetrics(ctx)
	if err != nil {
		slog.Error("Failed to initialize OpenTelemetry metrics", "error", err)
		// Continue without metrics if initialization fails
	} else {
		defer func() {
			// Give metrics time to flush on shutdown
			shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			if err := shutdownMetrics(shutdownCtx); err != nil {
				slog.Error("Failed to shutdown metrics provider", "error", err)
			} else {
				slog.Info("Metrics provider shutdown successfully")
			}
		}()
	}

	// Initialize traces
	shutdownTraces, err := otel.InitTraces(ctx)
	if err != nil {
		slog.Error("Failed to initialize OpenTelemetry traces", "error", err)
		// Continue without traces if initialization fails
	} else {
		defer func() {
			// Give traces time to flush on shutdown
			shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			if err := shutdownTraces(shutdownCtx); err != nil {
				slog.Error("Failed to shutdown trace provider", "error", err)
			} else {
				slog.Info("Trace provider shutdown successfully")
			}
		}()
	}

	// Setup graceful shutdown
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)

	// Run server in a goroutine
	errCh := make(chan error, 1)
	go func() {
		slog.Info("Starting server", "addr", app.Server.Addr)
		if err := app.Run(); err != nil {
			errCh <- err
		}
	}()

	// Wait for shutdown signal or error
	select {
	case err := <-errCh:
		slog.Error("Server error", "error", err)
		os.Exit(1)
	case sig := <-sigCh:
		slog.Info("Received shutdown signal", "signal", sig)
	}
}
