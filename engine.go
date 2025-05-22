package fuego

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/getkin/kin-openapi/openapi3gen"
)

// An EngineOption represents configures the behavior of an [Engine].
type EngineOption func(*Engine)

// NewEngine creates a new Engine with the given options.
// For example:
//
//	engine := fuego.NewEngine(
//		WithOpenAPIConfig(
//			OpenAPIConfig{
//				PrettyFormatJSON: true,
//			},
//		),
//	)
//
// Options all begin with `With`.
func NewEngine(options ...EngineOption) *Engine {
	e := &Engine{
		OpenAPI:              NewOpenAPI(),
		ErrorHandler:         ErrorHandler,
		responseContentTypes: defaultResponseContentTypes,
	}
	for _, option := range options {
		option(e)
	}
	return e
}

// The Engine is the main struct of the framework.
type Engine struct {
	OpenAPI      *OpenAPI
	ErrorHandler func(context.Context, error) error

	requestContentTypes  []string
	responseContentTypes []string
}

type OpenAPIConfig struct {
	// Local path to save the OpenAPI JSON spec
	JSONFilePath string
	// If true, the server will not serve nor generate any OpenAPI resources
	Disabled bool
	// If true, the engine will not print messages
	DisableMessages bool
	// If true, the engine will not save the OpenAPI JSON spec locally
	DisableLocalSave bool
	// If true, no default server will be added.
	// Note: this option only applies to the fuego [Server]. Adaptors are not affected by this option.
	DisableDefaultServer bool
	// Pretty prints the OpenAPI spec with proper JSON indentation
	PrettyFormatJSON bool
	// URL to serve the OpenAPI JSON spec
	SpecURL string
	// Handler to serve the OpenAPI UI from spec URL
	UIHandler func(specURL string) http.Handler
	// URL to serve the swagger UI
	SwaggerURL string
	// If true, the server will not serve the Swagger UI
	DisableSwaggerUI bool
	// Middleware configuration for the engine
	MiddlewareConfig MiddlewareConfig
}

var (
	defaultOpenAPIConfig = OpenAPIConfig{
		JSONFilePath: "doc/openapi.json",
		SpecURL:      "/swagger/openapi.json",
		SwaggerURL:   "/swagger",
		UIHandler:    DefaultOpenAPIHandler,
		MiddlewareConfig: MiddlewareConfig{
			DisableMiddlewareSection: false,
			MaxNumberOfMiddlewares:   6,
			ShortMiddlewaresPaths:    false,
		},
	}
	defaultResponseContentTypes = []string{"application/json", "application/xml"}
)

// WithRequestContentType sets the accepted content types for the engine.
// By default, the accepted content types is */*.
func WithRequestContentType(consumes ...string) EngineOption {
	return func(e *Engine) { e.requestContentTypes = consumes }
}

// WithResponseContentType sets content types of the returned body.
// By default, the returned content-types' are application/json and application/xml
func WithResponseContentType(consumes ...string) func(*Engine) {
	return func(e *Engine) { e.responseContentTypes = consumes }
}

type MiddlewareConfig struct {
	DisableMiddlewareSection bool
	MaxNumberOfMiddlewares   int
	ShortMiddlewaresPaths    bool
}

func WithMiddlewareConfig(cfg MiddlewareConfig) EngineOption {
	return func(e *Engine) {
		e.OpenAPI.Config.MiddlewareConfig.DisableMiddlewareSection = cfg.DisableMiddlewareSection
		e.OpenAPI.Config.MiddlewareConfig.ShortMiddlewaresPaths = cfg.ShortMiddlewaresPaths
		if cfg.MaxNumberOfMiddlewares != 0 {
			e.OpenAPI.Config.MiddlewareConfig.MaxNumberOfMiddlewares = cfg.MaxNumberOfMiddlewares
		}
	}
}

func WithOpenAPIConfig(config OpenAPIConfig) EngineOption {
	return func(e *Engine) {
		if config.JSONFilePath != "" {
			e.OpenAPI.Config.JSONFilePath = config.JSONFilePath
		}
		if config.SpecURL != "" {
			e.OpenAPI.Config.SpecURL = config.SpecURL
		}
		if config.SwaggerURL != "" {
			e.OpenAPI.Config.SwaggerURL = config.SwaggerURL
		}
		if config.UIHandler != nil {
			e.OpenAPI.Config.UIHandler = config.UIHandler
		}

		e.OpenAPI.Config.Disabled = config.Disabled
		e.OpenAPI.Config.DisableLocalSave = config.DisableLocalSave
		e.OpenAPI.Config.DisableDefaultServer = config.DisableDefaultServer
		e.OpenAPI.Config.PrettyFormatJSON = config.PrettyFormatJSON
		e.OpenAPI.Config.DisableSwaggerUI = config.DisableSwaggerUI
		e.OpenAPI.Config.DisableMessages = config.DisableMessages

		if !validateSpecURL(e.OpenAPI.Config.SpecURL) {
			slog.Error("Error serving OpenAPI JSON spec. Value of 's.OpenAPIServerConfig.SpecURL' option is not valid", "url", e.OpenAPI.Config.SpecURL)
			return
		}
		if !validateSwaggerURL(e.OpenAPI.Config.SwaggerURL) {
			slog.Error("Error serving Swagger UI. Value of 's.OpenAPIServerConfig.SwaggerURL' option is not valid", "url", e.OpenAPI.Config.SwaggerURL)
			return
		}

		WithMiddlewareConfig(config.MiddlewareConfig)(e)
	}
}

// WithOpenAPIGeneratorConfig sets the options for the generator.
func WithOpenAPIGeneratorSchemaCustomizer(sc openapi3gen.SchemaCustomizerFn) EngineOption {
	return func(e *Engine) {
		e.OpenAPI.SetGeneratorSchemaCustomizer(sc)
	}
}

// WithErrorHandler sets a customer error handler for the server
func WithErrorHandler(errorHandler func(ctx context.Context, err error) error) EngineOption {
	return func(e *Engine) {
		if errorHandler == nil {
			panic("errorHandler cannot be nil")
		}

		e.ErrorHandler = errorHandler
	}
}

// DisableErrorHandler overrides ErrorHandler with a simple pass-through
func DisableErrorHandler() EngineOption {
	return func(e *Engine) {
		e.ErrorHandler = func(_ context.Context, err error) error { return err }
	}
}

func (e *Engine) SpecHandler() func(c ContextNoBody) (openapi3.T, error) {
	return func(c ContextNoBody) (openapi3.T, error) {
		return *e.OpenAPI.Description(), nil
	}
}

// OutputOpenAPISpec takes the OpenAPI spec and outputs it to a JSON file
func (e *Engine) OutputOpenAPISpec() *openapi3.T {
	e.OpenAPI.computeTags()

	// Validate
	err := e.OpenAPI.Description().Validate(context.Background())
	if err != nil {
		slog.Error("Error validating spec", "error", err)
	}

	// Marshal spec to JSON
	jsonSpec, err := e.marshalSpec()
	if err != nil {
		slog.Error("Error marshaling spec to JSON", "error", err)
	}

	if !e.OpenAPI.Config.DisableLocalSave {
		err := e.saveOpenAPIToFile(e.OpenAPI.Config.JSONFilePath, jsonSpec)
		if err != nil {
			slog.Error("Error saving spec to local path", "error", err, "path", e.OpenAPI.Config.JSONFilePath)
		}
	}
	return e.OpenAPI.Description()
}

func (e *Engine) saveOpenAPIToFile(jsonSpecLocalPath string, jsonSpec []byte) error {
	jsonFolder := filepath.Dir(jsonSpecLocalPath)

	err := os.MkdirAll(jsonFolder, 0o750)
	if err != nil {
		return fmt.Errorf("error creating docs directory: %w", err)
	}

	f, err := os.Create(jsonSpecLocalPath) // #nosec G304 (file path provided by developer, not by user)
	if err != nil {
		return fmt.Errorf("error creating file: %w", err)
	}
	defer f.Close()

	_, err = f.Write(jsonSpec)
	if err != nil {
		return fmt.Errorf("error writing file: %w", err)
	}

	e.printOpenAPIMessage("JSON file: " + jsonSpecLocalPath)
	return nil
}

func (e *Engine) marshalSpec() ([]byte, error) {
	if e.OpenAPI.Config.PrettyFormatJSON {
		return json.MarshalIndent(e.OpenAPI.Description(), "", "\t")
	}
	return json.Marshal(e.OpenAPI.Description())
}

func (e *Engine) printOpenAPIMessage(msg string) {
	if !e.OpenAPI.Config.DisableMessages {
		slog.Info(msg)
	}
}
