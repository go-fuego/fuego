package lib

import (
	"context"
	"os"
	"testing"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/stretchr/testify/require"
	"gotest.tools/v3/golden"

	"github.com/go-fuego/fuego"
)

func TestPetstoreOpenAPIGeneration(t *testing.T) {
	server := NewPetStoreServer(
		fuego.WithoutStartupMessages(),
		fuego.WithEngineOptions(
			fuego.WithOpenAPIConfig(fuego.OpenAPIConfig{
				JSONFilePath:     "testdata/doc/openapi.json",
				PrettyFormatJSON: true,
				Info: &openapi3.Info{
					Title:   "Pet Store",
					Version: "0.0.2",
				},
			}),
		),
	)

	server.Engine.RegisterOpenAPIRoutes(server)
	server.OutputOpenAPISpec()
	err := server.OpenAPI.Description().Validate(context.Background())
	require.NoError(t, err)

	generatedSpec, err := os.ReadFile("testdata/doc/openapi.json")
	require.NoError(t, err)

	golden.Assert(t, string(generatedSpec), "doc/openapi.golden.json")
}
