package fuego

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"os"
	"path/filepath"
)

func NewEngine(options ...func(*Engine)) *Engine {
	e := &Engine{
		OpenAPI:       NewOpenAPI(),
		OpenAPIConfig: defaultOpenAPIConfig,
		ErrorHandler:  ErrorHandler,
	}
	for _, option := range options {
		option(e)
	}
	return e
}

// The Engine is the main struct of the framework.
type Engine struct {
	OpenAPI       *OpenAPI
	ErrorHandler  func(error) error
	OpenAPIConfig OpenAPIConfig

	acceptedContentTypes []string
}

type EngineOpenAPIConfig struct {
	// If true, the engine will not print messages
	DisableMessages bool
	// If true, the engine will not save the OpenAPI JSON spec locally
	DisableLocalSave bool
	// Local path to save the OpenAPI JSON spec
	JSONFilePath string
	// Pretty prints the OpenAPI spec with proper JSON indentation
	PrettyFormatJSON bool
}

func WithOpenAPIConfig(config OpenAPIConfig) func(*Engine) {
	return func(e *Engine) {
		if config.JsonURL != "" {
			e.OpenAPIConfig.JsonURL = config.JsonURL
		}

		if config.SwaggerURL != "" {
			e.OpenAPIConfig.SwaggerURL = config.SwaggerURL
		}

		if config.JSONFilePath != "" {
			e.OpenAPIConfig.JSONFilePath = config.JSONFilePath
		}

		if config.UIHandler != nil {
			e.OpenAPIConfig.UIHandler = config.UIHandler
		}

		e.OpenAPIConfig.DisableSwagger = config.DisableSwagger
		e.OpenAPIConfig.DisableSwaggerUI = config.DisableSwaggerUI
		e.OpenAPIConfig.DisableLocalSave = config.DisableLocalSave
		e.OpenAPIConfig.PrettyFormatJSON = config.PrettyFormatJSON

		if !validateJsonSpecUrl(e.OpenAPIConfig.JsonURL) {
			slog.Error("Error serving openapi json spec. Value of 's.OpenAPIConfig.JsonSpecUrl' option is not valid", "url", e.OpenAPIConfig.JsonURL)
			return
		}

		if !validateSwaggerUrl(e.OpenAPIConfig.SwaggerURL) {
			slog.Error("Error serving swagger ui. Value of 's.OpenAPIConfig.SwaggerUrl' option is not valid", "url", e.OpenAPIConfig.SwaggerURL)
			return
		}
	}
}

// OutputOpenAPISpec takes the OpenAPI spec and outputs it to a JSON file
func (e *Engine) OutputOpenAPISpec() []byte {
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

	if !e.OpenAPIConfig.DisableLocalSave {
		err := e.saveOpenAPIToFile(e.OpenAPIConfig.JSONFilePath, jsonSpec)
		if err != nil {
			slog.Error("Error saving spec to local path", "error", err, "path", e.OpenAPIConfig.JSONFilePath)
		}
	}
	return jsonSpec
}

func (e *Engine) saveOpenAPIToFile(jsonSpecLocalPath string, jsonSpec []byte) error {
	jsonFolder := filepath.Dir(jsonSpecLocalPath)

	err := os.MkdirAll(jsonFolder, 0o750)
	if err != nil {
		return errors.New("error creating docs directory")
	}

	f, err := os.Create(jsonSpecLocalPath) // #nosec G304 (file path provided by developer, not by user)
	if err != nil {
		return errors.New("error creating file")
	}
	defer f.Close()

	_, err = f.Write(jsonSpec)
	if err != nil {
		return errors.New("error writing file ")
	}

	e.printOpenAPIMessage("JSON file: " + jsonSpecLocalPath)
	return nil
}

func (s *Engine) marshalSpec() ([]byte, error) {
	if s.OpenAPIConfig.PrettyFormatJSON {
		return json.MarshalIndent(s.OpenAPI.Description(), "", "	")
	}
	return json.Marshal(s.OpenAPI.Description())
}

func (e *Engine) printOpenAPIMessage(msg string) {
	if !e.OpenAPIConfig.DisableMessages {
		slog.Info(msg)
	}
}
