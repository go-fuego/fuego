package main

import (
	"github.com/go-fuego/fuego/examples/petstore/lib"
)

func main() {
	err := lib.NewPetStoreServer().Run()
	if err != nil {
		panic(err)
	}
}
