---
title: Hello World
---

# Hello world

How to write Hello World with Fuego.

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
