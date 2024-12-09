package lib

import (
	"github.com/gin-gonic/gin"

	"github.com/go-fuego/fuego"
	"github.com/go-fuego/fuego/extra/fuegogin"
)

type HelloRequest struct {
	Name string `json:"name"`
}

type HelloResponse struct {
	Message string `json:"message"`
}

func SetupGin() (*gin.Engine, *fuego.Server) {
	e := gin.Default()
	s := fuego.NewServer()

	e.GET("/gin", ginController)
	fuegogin.Get(s, e, "/fuego", fuegoController)

	return e, s
}

func ginController(c *gin.Context) {
	c.String(200, "pong")
}

func fuegoController(c *fuegogin.ContextWithBody[HelloRequest]) (HelloResponse, error) {
	return HelloResponse{
		Message: "Hello ",
	}, nil
}
