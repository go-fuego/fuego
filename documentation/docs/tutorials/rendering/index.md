# HTML Rendering

Fuego is not only capable to handle XML and JSON, it can also render HTML.

It supports templating with [html/template](https://pkg.go.dev/html/template),
[Templ](https://github.com/a-h/templ), and [Gomponents](https://github.com/maragudk/gomponents).

## Content Negotiation

Remember that Fuego handles [Content Negotiation](https://developer.mozilla.org/en-US/docs/Web/HTTP/Content_negotiation),
so you can serve different content types based on the `Accept` header of the request.

Fuego also provides a helper to render HTML or JSON data for a single controller!

```go
package main

import (
	"github.com/go-fuego/fuego"
)

func main() {
	s := fuego.NewServer()

	fuego.Get(s, "/", func(c fuego.ContextNoBody) (interface{}, error) {
		return fuego.DataOrHTML(
			data, // When asking for JSON/XML, this data will be returned
			MyTemplateInjectedWithData(data), // When asking for HTML, this template will be rendered
		), nil
	})

	s.Run()
}
```
