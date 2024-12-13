package main

import (
	"github.com/go-fuego/fuego/extra/fuegogin/lib"
)

func main() {
	e, _ := lib.SetupGin()

	err := e.Run(":8980")
	if err != nil {
		panic(err)
	}
}
