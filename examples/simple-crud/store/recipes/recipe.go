package recipes

import (
	"strings"

	"github.com/go-fuego/fuego"
	"github.com/google/uuid"
)

var _ fuego.InTransformer = (*CreateRecipeParams)(nil)

// InTransform implements fuego.InTransformer.
func (c *CreateRecipeParams) InTransform() error {
	c.ID = uuid.NewString()
	c.Name = strings.TrimSpace(c.Name)

	return nil
}
