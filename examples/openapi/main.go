package main

import (
	"github.com/go-fuego/fuego"
	"github.com/go-fuego/fuego/option"
)

func main() {
	s := fuego.NewServer()

	fuego.Get(s, "/", helloWorld,
		option.Summary("A simple hello world"),
		option.Description("This is a simple hello world"),
		option.Deprecated(),
	)

	s.Run()
}

func helloWorld(c fuego.ContextNoBody) (string, error) {
	return "Hello, World!", nil
}
