package handler

import (
	"context"

	"github.com/go-fuego/fuego"
)

// usernameFromContext is a helper function that extracts the username from the context.
func usernameFromContext(ctx context.Context) (string, error) {
	t, err := fuego.TokenFromContext(ctx)
	if err != nil {
		return "", err
	}

	return t.GetIssuer()
}
