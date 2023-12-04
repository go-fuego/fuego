package store

import (
	"strings"

	"github.com/go-fuego/fuego"
)

var _ fuego.InTransformer = (*CreateIngredientParams)(nil)

func (c *CreateIngredientParams) InTransform() error {
	c.Name = strings.TrimSpace(c.Name)

	c.ID = slug(c.Name)

	return nil
}
