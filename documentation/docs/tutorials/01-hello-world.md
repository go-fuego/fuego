# Hello world

Let's discover **Fuego** in a few lines.

## Quick start

If you don't want to copy/paste the code on your local setup, you can run the
following command:

```bash
go run github.com/go-fuego/fuego/examples/hello-world@latest
```

Useful URLs (including OpenAPI spec & Swagger UI) are given in the terminal:
you'll be able to see the result in your browser.

## Start from scratch

First, create a directory for your project:

```bash
mkdir hello-fuego
cd hello-fuego
```

Then, create a `go.mod` file:

```bash

go mod init hello-fuego
```

Finally, create a `main.go` file with the following content:

```go
package main

import (
	"github.com/go-fuego/fuego"
)

func main() {
	s := fuego.NewServer()

	fuego.Get(s, "/", func(c fuego.ContextNoBody) (string, error) {
		return "Hello, World!", nil
	})

	s.Run()
}
```

You can now run your server:

```bash
go mod tidy
go run .
```
