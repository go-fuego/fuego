package main

import (
	"fmt"

	"github.com/go-fuego/fuego"
)

func main() {
	fmt.Println("curl -X GET http://localhost:9999/")
	fmt.Println("----------------------------------")

	s := fuego.NewServer()

	fuego.Get(s, "/", func(c fuego.Ctx[any]) (string, error) {
		return "Hello, World!", nil
	})

	s.Run()
}
