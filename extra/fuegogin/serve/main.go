package main

import "github.com/go-fuego/fuego/extra/fuegogin/lib"

func main() {
	e, s := lib.SetupGin()

	go func() {
		s.Run()
	}()

	err := e.Run(":8080")
	if err != nil {
		panic(err)
	}
}
