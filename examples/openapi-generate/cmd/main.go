package main

import (
	"github.com/go-fuego/fuego/examples/openapi-generate/server"
)

func main() {
	s := server.GetServer()
	s.Run()
}
