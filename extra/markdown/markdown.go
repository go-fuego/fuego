package markdown

import (
	"html/template"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
)

// Markdown converts a markdown string to HTML.
// Note: fuego does not protect against malicious content
// sanitation is up the caller of this function.
func Markdown(content string) template.HTML {
	if content == "" {
		return template.HTML("")
	}
	mdRenderer := html.NewRenderer(html.RendererOptions{Flags: html.CommonFlags | html.SkipHTML})
	mdParser := parser.NewWithExtensions(parser.CommonExtensions | parser.AutoHeadingIDs | parser.NoEmptyLineBeforeBlock | parser.Footnotes | parser.DefinitionLists)

	//nolint:gosec // G203 // the caller of this function needs to sanitize their input
	return template.HTML(markdown.ToHTML([]byte(content), mdParser, mdRenderer))
}
