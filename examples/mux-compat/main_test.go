package main

import (
	"testing"

	"github.com/go-fuego/fuego"
	"github.com/stretchr/testify/require"
)

func TestFuegoControllerPost(t *testing.T) {
	testCtx := fuego.NewMockContext(HelloRequest{Word: "World"}, any(nil))
	testCtx.QueryParams().Set("name", "Ewen")

	response, err := fuegoControllerPost(testCtx)
	require.NoError(t, err)
	require.Equal(t, "Hello world, Ewen", response.Message)
}
