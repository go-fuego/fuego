package main

import (
	"fmt"

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

func main() {
	e, _ := server()

	fmt.Println("OpenAPI at at http://localhost:8980/swagger ✅")

	err := e.Run(":8980")
	if err != nil {
		panic(err)
	}
}

func server() (*gin.Engine, *fuego.OpenAPI) {
	ginRouter := gin.Default()
	engine := fuego.NewEngine()

	// Register Gin controller
	ginRouter.GET("/gin", ginController)

	// Incrementally add OpenAPI spec
	// Level 1: Register Gin controller to Gin router, plugs Fuego OpenAPI route declaration
	fuegogin.GetGin(engine, ginRouter, "/gin-with-openapi", ginController)

	// Level 2: Register Fuego controller to Gin router. Fuego take care of serialization/deserialization, error handling, content-negotiation, etc.
	fuegogin.Get(engine, ginRouter, "/fuego", fuegoControllerGet)

	// Add some options to the POST endpoint
	fuegogin.Post(engine, ginRouter, "/fuego-with-options", fuegoControllerPost,
		// OpenAPI options
		option.Description("Some description"),
		option.OperationID("SomeOperationID"),
		option.AddError(409, "Name Already Exists"),
		option.DefaultStatusCode(201),
		option.Tags("Fuego"),

		// Add some parameters.
		option.Query("name", "Your name", param.Example("name example", "John Carmack")),
		option.Header("X-Request-ID", "Request ID", param.Default("123456")),
		option.Header("Content-Type", "Content Type", param.Default("application/json")),
	)

	// Supports groups & path parameters even for gin handlers
	group := ginRouter.Group("/my-group/:id")
	fuegogin.Get(engine, group, "/fuego", fuegoControllerGet,
		option.Summary("Route with group and id"),
	)

	// Serve the OpenAPI spec
	ginRouter.GET("/openapi.json", serveOpenApiJSONDescription(engine.OpenAPI))
	ginRouter.GET("/swagger", DefaultOpenAPIHandler("/openapi.json"))

	return ginRouter, engine.OpenAPI
}
