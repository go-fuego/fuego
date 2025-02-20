package main

import (
	"github.com/go-fuego/fuego/examples/openapi-generate/server"
)

func main() {
	// Get the server instance, configured earlier
	newServer := server.GetServer()

	// Simple configuration for OpenAPI spec generation
	newServer.OpenAPI.Config.DisableLocalSave = false
	newServer.OpenAPI.Config.PrettyFormatJSON = true
	newServer.OpenAPI.Config.JSONFilePath = "api/openapi.json"

	// Generate the OpenAPI spec
	newServer.OutputOpenAPISpec()
}
