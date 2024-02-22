package main

import (
	"github.com/go-fuego/fuego"
)

func main() {
	s := fuego.NewServer()

	fuego.Get(s, "/", helloWorld).
		Summary("A simple hello world").
		Description("This is a simple hello world").
		Deprecated()

	s.Run()
}

func helloWorld(c fuego.ContextNoBody) (string, error) {
	return "Hello, World!", nil
}
