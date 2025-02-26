# Routing

## `net/http` (default)

Fuego uses the `fuego.Get`, `fuego.Post`, `fuego.Put`, `fuego.Patch`, `fuego.Delete`, `fuego.Options`, `fuego.Head` and `fuego.Any` functions to define routes. Example:

```go
package main

import (
	"github.com/go-fuego/fuego"
)

func main() {
	s := fuego.NewServer()

	fuego.Get(s, "/books", helloWorld)
	fuego.Post(s, "/books", helloWorld)
	fuego.Get(s, "/books/{id}", helloWorld)

	s.Run()
}
```
