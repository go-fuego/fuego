package fuego

import (
	"html/template"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
)

var mdRenderer = html.NewRenderer(html.RendererOptions{Flags: html.CommonFlags | html.SkipHTML})

// Markdown converts a markdown string to HTML.
func Markdown(content string) template.HTML {
	if content == "" {
		return template.HTML("")
	}
	mdParser := parser.NewWithExtensions(parser.CommonExtensions | parser.AutoHeadingIDs | parser.NoEmptyLineBeforeBlock | parser.Footnotes | parser.DefinitionLists)

	return template.HTML(markdown.ToHTML([]byte(content), mdParser, mdRenderer))
}
