package store

import (
	"strings"

	"github.com/go-op/op"
	"github.com/google/uuid"
)

var _ op.InTransformer = (*CreateRecipeParams)(nil)

// InTransform implements op.InTransformer.
func (c *CreateRecipeParams) InTransform() error {
	c.ID = uuid.NewString()
	c.Name = strings.TrimSpace(c.Name)

	return nil
}
