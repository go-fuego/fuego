package ingredients

import (
	"strings"

	"github.com/go-fuego/fuego"
	"github.com/google/uuid"
)

var _ fuego.InTransformer = (*CreateIngredientParams)(nil)

func (c *CreateIngredientParams) InTransform() error {
	c.ID = uuid.NewString()
	c.Name = strings.TrimSpace(c.Name)

	return nil
}
