package fuego

import (
	"html/template"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMarkdown(t *testing.T) {
	t.Run("empty string", func(t *testing.T) {
		html := Markdown("")
		require.Equal(t, template.HTML(""), html)
	})

	t.Run("can render markdown", func(t *testing.T) {
		md := `# Hello
	Just **testing**.`

		html := Markdown(md)
		require.Equal(t, template.HTML("<h1 id=\"hello\">Hello</h1>\n\n<pre><code>Just **testing**.\n</code></pre>\n"), html)
	})
}
