package main

import (
	"fmt"
	"net/http"

	"github.com/go-fuego/fuego"
)

func muxController(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "pong")
}

func fuegoControllerGet(c fuego.ContextNoBody) (HelloResponse, error) {
	return HelloResponse{
		Message: "Hello",
	}, nil
}

func fuegoControllerPost(c fuego.ContextWithBody[HelloRequest]) (*HelloResponse, error) {
	body, err := c.Body()
	if err != nil {
		return nil, err
	}

	if body.Word == "forbidden" {
		return nil, fuego.BadRequestError{Title: "Forbidden word"}
	}

	name := c.QueryParam("name")

	return &HelloResponse{
		Message: fmt.Sprintf("Hello %s, %s", body.Word, name),
	}, nil
}
