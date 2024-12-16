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
	Word string `json:"word"`
}

type HelloResponse struct {
	Message string `json:"message"`
}

func SetupGin() (*gin.Engine, *fuego.OpenAPI) {
	e := gin.Default()
	openapi := fuego.NewOpenAPI()

	// Register Gin controller
	e.GET("/gin", ginController)

	group := e.Group("/my-group/:id")
	fuegogin.Get(openapi, group, "/fuego", fuegoControllerGet)

	// Register to Gin router with Fuego wrapper for same OpenAPI spec
	fuegogin.Get(openapi, e, "/fuego", fuegoControllerGet)
	fuegogin.Post(openapi, e, "/fuego-with-options", fuegoControllerPost,
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

	fmt.Println("OpenAPI at at http://localhost:8980/swagger")

	return e, openapi
}

func ginController(c *gin.Context) {
	c.String(200, "pong")
}

func fuegoControllerGet(c fuegogin.ContextNoBody) (HelloResponse, error) {
	return HelloResponse{
		Message: "Hello",
	}, nil
}

func fuegoControllerPost(c fuegogin.ContextWithBody[HelloRequest]) (HelloResponse, error) {
	body, err := c.Body()
	if err != nil {
		return HelloResponse{}, err
	}

	name := c.QueryParam("name")

	return HelloResponse{
		Message: fmt.Sprintf("Hello %s, %s", body.Word, name),
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
