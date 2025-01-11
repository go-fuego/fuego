package main

import (
	"fmt"

	"github.com/gin-gonic/gin"

	"github.com/go-fuego/fuego"
)

func ginController(c *gin.Context) {
	c.String(200, "pong")
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

	_ = c.Context().(*gin.Context) // Access to the Gin context

	name := c.QueryParam("name")
	_ = c.QueryParam("not-existing-param-raises-warning")

	return &HelloResponse{
		Message: fmt.Sprintf("Hello %s, %s", body.Word, name),
	}, nil
}
