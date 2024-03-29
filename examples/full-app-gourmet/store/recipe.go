package store

import (
	"context"
	"strings"

	"github.com/go-fuego/fuego"
)

var _ fuego.InTransformer = (*CreateRecipeParams)(nil)

// InTransform implements fuego.InTransformer.
func (c *CreateRecipeParams) InTransform(context.Context) error {
	c.Name = strings.TrimSpace(c.Name)

	c.ID = slug(c.Name)

	return nil
}
