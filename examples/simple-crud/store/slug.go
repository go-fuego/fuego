package store

import (
	"net/url"
	"strings"
	"unicode"

	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

var transformer = transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)), norm.NFC)

// slug returns a slugified version of name.
func slug(name string) string {
	name = strings.TrimSpace(name)

	id := strings.ToLower(name)
	id = strings.ReplaceAll(id, " ", "-")

	var err error
	id, _, err = transform.String(transformer, id)
	if err != nil {
		panic(err)
	}
	id = url.PathEscape(id)

	return id
}
