package main

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/go-op/op"
)

type bod struct {
	Name string `json:"name" validate:"required"`
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

func stdController(w http.ResponseWriter, r *http.Request) {
	slog.Info("controller")

	message := "Hello World."
	limit, ok := r.URL.Query()["limit"]
	if ok {
		message += " The limit is:" + limit[0]
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(ans{Ans: message})
}

func main() {
	s := op.NewServer(
		op.WithPort(":8070"),
		op.WithDisallowUnknownFields(false),
	)

	op.Get(s, "/hello", controller)
	op.Post(s, "/hello", controller)
	op.Post(s, "/hello2", controller2)

	s.Run()
}
