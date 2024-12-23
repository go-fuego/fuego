package fuego

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"os"
	"path/filepath"
)

func NewEngine(config ...OpenAPIConfig) *Engine {
	if len(config) > 1 {
		panic("config should not be more than one")
	}
	engine := &Engine{
		OpenAPI:       NewOpenAPI(),
		OpenAPIConfig: defaultOpenAPIConfig,
		ErrorHandler:  ErrorHandler,
	}
	if len(config) > 0 {
		engine.setOpenAPIConfig(config[0])
	}
	return engine
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
	JsonFilePath string
	// Pretty prints the OpenAPI spec with proper JSON indentation
	PrettyFormatJson bool
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
		err := e.saveOpenAPIToFile(e.OpenAPIConfig.JsonFilePath, jsonSpec)
		if err != nil {
			slog.Error("Error saving spec to local path", "error", err, "path", e.OpenAPIConfig.JsonFilePath)
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
	if s.OpenAPIConfig.PrettyFormatJson {
		return json.MarshalIndent(s.OpenAPI.Description(), "", "	")
	}
	return json.Marshal(s.OpenAPI.Description())
}

func (e *Engine) setOpenAPIConfig(config OpenAPIConfig) {
	if config.JsonUrl != "" {
		e.OpenAPIConfig.JsonUrl = config.JsonUrl
	}

	if config.SwaggerUrl != "" {
		e.OpenAPIConfig.SwaggerUrl = config.SwaggerUrl
	}

	if config.JsonFilePath != "" {
		e.OpenAPIConfig.JsonFilePath = config.JsonFilePath
	}

	if config.UIHandler != nil {
		e.OpenAPIConfig.UIHandler = config.UIHandler
	}

	e.OpenAPIConfig.DisableSwagger = config.DisableSwagger
	e.OpenAPIConfig.DisableSwaggerUI = config.DisableSwaggerUI
	e.OpenAPIConfig.DisableLocalSave = config.DisableLocalSave
	e.OpenAPIConfig.PrettyFormatJson = config.PrettyFormatJson

	if !validateJsonSpecUrl(e.OpenAPIConfig.JsonUrl) {
		slog.Error("Error serving openapi json spec. Value of 's.OpenAPIConfig.JsonSpecUrl' option is not valid", "url", e.OpenAPIConfig.JsonUrl)
		return
	}

	if !validateSwaggerUrl(e.OpenAPIConfig.SwaggerUrl) {
		slog.Error("Error serving swagger ui. Value of 's.OpenAPIConfig.SwaggerUrl' option is not valid", "url", e.OpenAPIConfig.SwaggerUrl)
		return
	}
}

func (e *Engine) printOpenAPIMessage(msg string) {
	if !e.OpenAPIConfig.DisableMessages {
		slog.Info(msg)
	}
}
