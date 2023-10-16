package main

import (
	"log/slog"
	"op/op"
)

type bod struct {
	Name string `json:"name"`
}

type ans struct {
	Ans string `json:"ans"`
}

func controller(c op.Ctx[any]) (ans, error) {
	slog.Info("controller")

	message := "Hello World."
	limit, ok := c.QueryParams()["limit"]
	if ok {
		message += " The limit is:" + limit
	}
	return ans{Ans: message}, nil
}

func controllerPost(c op.Ctx[bod]) (ans, error) {
	slog.Info("controller")
	return ans{Ans: "Wooow " + c.Body().Name}, nil
}

func controller2(c op.Ctx[bod]) (string, error) {

	return "Hello " + c.Body().Name, nil
}

func main() {

	op.Get("/hello", controller)
	// op.Post("/hello", controllerPost)

	op.Post("/hello2", controller2)

	op.Run(":8080")
}
