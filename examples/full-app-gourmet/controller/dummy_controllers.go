package controller

import (
	"time"

	"github.com/go-fuego/fuego"
)

type test struct {
	Name string `json:"name"`
}

func slow(c fuego.ContextNoBody) (test, error) {
	time.Sleep(2 * time.Second)
	return test{Name: "hello"}, nil
}

func placeholderController(c fuego.ContextNoBody) (string, error) {
	return "hello", nil
}
