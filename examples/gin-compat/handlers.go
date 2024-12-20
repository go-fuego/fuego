package main

import (
	"fmt"
	"net/http"

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

func fuegoControllerPost(c fuego.ContextWithBody[HelloRequest]) (HelloResponse, error) {
	body, err := c.Body()
	if err != nil {
		return HelloResponse{}, err
	}

	if body.Word == "forbidden" {
		return HelloResponse{}, fuego.BadRequestError{Title: "Forbidden word"}
	}

	ctx := c.Context().(*gin.Context)
	fmt.Printf("%#v", ctx)

	name := c.QueryParam("name")

	return HelloResponse{
		Message: fmt.Sprintf("Hello %s, %s", body.Word, name),
	}, nil
}

func serveOpenApiJSONDescription(s *fuego.OpenAPI) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, s.Description())
	}
}

func DefaultOpenAPIHandler(specURL string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Header("Content-Type", "text/html; charset=utf-8")
		ctx.String(200, fuego.DefaultOpenAPIHTML(specURL))
	}
}
