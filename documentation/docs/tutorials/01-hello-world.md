---
title: Hello World
position: 1
order: 1
---

# Hello world

Let's discover **Fuego** in a few lines.

## Run

If you don't want to copy/paste the code on your local setup, you can run the following command:

```bash
go run github.com/go-fuego/fuego/examples/hello-world@latest
```

Useful URLs are given in the terminal, you'll be able to see the result in your browser.

## Code

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
