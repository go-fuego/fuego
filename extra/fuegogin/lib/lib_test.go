package lib

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/go-fuego/fuego/extra/fuegogin"
)

func TestFuegoControllerPost(t *testing.T) {
	testCtx := &fuegogin.ContextTest[HelloRequest]{
		BodyInjected: HelloRequest{Word: "World"},
		Params:       url.Values{"name": []string{"Ewen"}},
	}

	response, err := fuegoControllerPost(testCtx)
	require.NoError(t, err)
	require.Equal(t, "Hello World, Ewen", response.Message)
}
