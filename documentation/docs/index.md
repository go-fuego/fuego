---
sidebar_position: 1
---

# 🔥 Fuego

<p align="center">
  <img src="/fuego/img/logo.svg" height="200" alt="Fuego Logo" />
</p>

Let's discover **Fuego in less than 5 minutes**.

## Quick peek without installing

Try our [Hello World](./tutorials/01-hello-world.md)!

```bash
go run github.com/go-fuego/fuego/examples/hello-world@latest
```

This runs the code for a simple hello world server.
Look at all it generates from a simple code!
You'll get a URL to see the result in your browser.

```go showLineNumbers
package main

import (
	"github.com/go-fuego/fuego"
)

func main() {
	s := fuego.NewServer()

	fuego.Get(s, "/", helloWorld)

	s.Run()
}

func helloWorld(c fuego.ContextNoBody) (string, error) {
	return "Hello, World!", nil
}
```

## Try example from real Fuego source code in 3 sec

Try Fuego immediately by cloning [the repo](https://github.com/go-fuego/fuego)
and running one of our examples.

```bash
git clone git@github.com:go-fuego/fuego.git
cd fuego/examples/petstore
go run .
```

### What you'll need

- [Golang v1.22](https://golang.org/doc/go1.22) or above
  _(Fuego relies on a new feature of the net/http package only available after 1.22)_.
