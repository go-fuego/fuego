package main

import (
	"log/slog"
	"op/op"
)

type bod struct {
	Name string `json:"name"`
}

func (b *bod) Normalize() error {
	b.Name = "normalized " + b.Name
	return nil
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

	me, err := c.Body()
	if err != nil {
		return ans{}, err
	}

	return ans{Ans: "Wooow " + me.Name}, nil
}

func controller2(c op.Ctx[bod]) (string, error) {
	me := c.MustBody()

	return "Hello " + me.Name, nil
}

func main() {

	op.Options(op.WithDisallowUnknownFields(false))

	op.Get("/hello", controller)
	// op.Post("/hello", controllerPost)

	op.Post("/hello2", controllerPost)

	op.Run(":8080")
}
