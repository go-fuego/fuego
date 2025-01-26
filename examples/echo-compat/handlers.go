package main

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/go-fuego/fuego"
)

func echoController(c echo.Context) error {
	return c.String(http.StatusOK, "pong")
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

func serveOpenApiJSONDescription(s *fuego.OpenAPI) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		return ctx.JSON(http.StatusOK, s.Description())
	}
}

func DefaultOpenAPIHandler(specURL string) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		ctx.Response().Header().Set(echo.HeaderContentType, "text/html; charset=utf-8")
		return ctx.String(http.StatusOK, fuego.DefaultOpenAPIHTML(specURL))
	}
}