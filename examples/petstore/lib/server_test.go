package lib

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"gotest.tools/v3/golden"

	"github.com/go-fuego/fuego"
)

func TestPetstoreOpenAPIGeneration(t *testing.T) {
	server := NewPetStoreServer(
		fuego.WithoutStartupMessages(),
		fuego.WithOpenAPIConfig(fuego.OpenAPIConfig{
			JsonFilePath:     "testdata/doc/openapi.json",
			PrettyFormatJson: true,
		}),
	)

	server.OutputOpenAPISpec()
	err := server.OpenAPIzer.OpenAPIDescription().Validate(context.Background())
	require.NoError(t, err)

	generatedSpec, err := os.ReadFile("testdata/doc/openapi.json")
	require.NoError(t, err)

	golden.Assert(t, string(generatedSpec), "doc/openapi.golden.json")
}
