package store

import (
	"net/url"
	"strings"
	"unicode"

	"github.com/go-fuego/fuego"
	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

var _ fuego.InTransformer = (*CreateRecipeParams)(nil)

// InTransform implements fuego.InTransformer.
func (c *CreateRecipeParams) InTransform() error {
	c.Name = strings.TrimSpace(c.Name)

	c.ID = strings.ToLower(c.Name)
	c.ID = strings.ReplaceAll(c.ID, " ", "-")

	t := transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)), norm.NFC)
	var err error
	c.ID, _, err = transform.String(t, c.ID)
	if err != nil {
		return err
	}
	c.ID = url.PathEscape(c.ID)

	return nil
}
