package lib

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/go-fuego/fuego"
	"github.com/go-fuego/fuego/extra/fuegogin"
	"github.com/go-fuego/fuego/option"
	"github.com/go-fuego/fuego/param"
)

type HelloRequest struct {
	Name string `json:"name"`
}

type HelloResponse struct {
	Message string `json:"message"`
}

func SetupGin() (*gin.Engine, *fuego.OpenAPI) {
	e := gin.Default()
	openapi := fuego.NewOpenAPI()

	// Register Gin controller
	e.GET("/gin", ginController)

	// Register to Gin router with Fuego wrapper for same OpenAPI spec
	fuegogin.Get(openapi, e, "/fuego", fuegoController)
	fuegogin.Get(openapi, e, "/fuego-with-options", fuegoController,
		option.Description("Some description"),
		option.OperationID("SomeOperationID"),
		option.AddError(409, "Name Already Exists"),
		option.DefaultStatusCode(201),
		option.Query("name", "Your name", param.Example("name example", "John Carmack")),
		option.Header("X-Request-ID", "Request ID", param.Default("123456")),
		option.Header("Content-Type", "Content Type", param.Default("application/json")),
	)

	// Serve the OpenAPI spec
	e.GET("/openapi.json", serveController(openapi))
	e.GET("/swagger", DefaultOpenAPIHandler("/openapi.json"))

	return e, openapi
}

func ginController(c *gin.Context) {
	c.String(200, "pong")
}

func fuegoController(c *fuegogin.ContextWithBody[HelloRequest]) (HelloResponse, error) {
	body, err := c.Body()
	if err != nil {
		return HelloResponse{}, err
	}
	fmt.Println("body", body)

	name := c.QueryParam("name")
	fmt.Println("name", name)

	return HelloResponse{
		Message: "Hello " + body.Name,
	}, nil
}

func serveController(s *fuego.OpenAPI) func(ctx *gin.Context) {
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
