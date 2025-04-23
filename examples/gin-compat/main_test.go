package main

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/go-fuego/fuego"
)

func TestFuegoControllerPost(t *testing.T) {
	testCtx := fuego.NewMockContext(HelloRequest{Word: "World"}, any(nil))
	testCtx.SetQueryParam("name", "Ewen")

	response, err := fuegoControllerPost(testCtx)
	require.NoError(t, err)
	require.Equal(t, "Hello World, Ewen", response.Message)
}
