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
	e := gin.Default()
	openapi := fuego.NewOpenAPI()

	// Register Gin controller
	e.GET("/gin", ginController)

	// Incrementally add OpenAPI spec
	// Level 1: Register Gin controller to Gin router, plugs Fuego OpenAPI route declaration
	fuegogin.GetGin(openapi, e, "/gin-with-openapi", ginController)

	// Level 2: Register Fuego controller to Gin router. Fuego take care of serialization/deserialization, error handling, content-negotiation, etc.
	fuegogin.Get(openapi, e, "/fuego", fuegoControllerGet)

	// Add some options to the POST endpoint
	fuegogin.Post(openapi, e, "/fuego-with-options", fuegoControllerPost,
		// OpenAPI options
		option.Description("Some description"),
		option.OperationID("SomeOperationID"),
		option.AddError(409, "Name Already Exists"),
		option.DefaultStatusCode(201),

		// Add some parameters.
		option.Query("name", "Your name", param.Example("name example", "John Carmack")),
		option.Header("X-Request-ID", "Request ID", param.Default("123456")),
		option.Header("Content-Type", "Content Type", param.Default("application/json")),
	)

	// Supports groups & path parameters even for gin handlers
	group := e.Group("/my-group/:id")
	fuegogin.Get(openapi, group, "/fuego", fuegoControllerGet,
		option.Summary("Route with group and id"),
	)

	// Serve the OpenAPI spec
	e.GET("/openapi.json", serveController(openapi))
	e.GET("/swagger", DefaultOpenAPIHandler("/openapi.json"))

	fmt.Println("OpenAPI at at http://localhost:8980/swagger ✅")

	err := e.Run(":8980")
	if err != nil {
		panic(err)
	}
}
