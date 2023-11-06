package controller

import (
	"time"

	"github.com/go-op/op"
)

type test struct {
	Name string `json:"name"`
}

func slow(c op.Ctx[any]) (test, error) {
	time.Sleep(2 * time.Second)
	return test{Name: "hello"}, nil
}

func placeholderController(c op.Ctx[any]) (string, error) {
	return "hello", nil
}
