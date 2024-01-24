package main

import (
	"fmt"

	"github.com/go-fuego/fuego"
)

func main() {
	fmt.Println("curl -X GET http://localhost:9999/")
	fmt.Println("----------------------------------")

	s := fuego.NewServer()

	fuego.Get(s, "/", func(c *fuego.ContextNoBody) (string, error) {
		return "Hello, World!", nil
	})

	s.Run()
}
