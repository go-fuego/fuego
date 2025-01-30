package main

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/go-fuego/fuego"
	"github.com/go-fuego/fuego/extra/fuegoecho"
	"github.com/go-fuego/fuego/option"
	"github.com/go-fuego/fuego/param"
)

type HelloRequest struct {
	Word string `json:"word" validate:"required,min=2"`
}

var _ fuego.InTransformer = &HelloRequest{}

type HelloResponse struct {
	Message string `json:"message"`
}

func main() {
	e, _ := server()

	fmt.Println("OpenAPI at http://localhost:8980/swagger ✅")

	err := e.Start(":8980")
	if err != nil {
		panic(err)
	}
}

func server() (*echo.Echo, *fuego.OpenAPI) {
	echoRouter := echo.New()
	echoRouter.Use(middleware.Logger())
	echoRouter.Use(middleware.Recover())

	engine := fuego.NewEngine()

	// Register Echo controller
	echoRouter.GET("/echo", echoController)

	// Incrementally add OpenAPI spec
	// 1️⃣ Level 1: Register Echo controller to Echo router, plugs Fuego OpenAPI route declaration
	fuegoecho.GetEcho(engine, echoRouter, "/echo-with-openapi", echoController)

	// 2️⃣ Level 2: Register Echo controller to Echo router, manually add options (not checked inside the Echo controller)
	fuegoecho.GetEcho(engine, echoRouter, "/echo-with-openapi-and-options", echoController,
		// OpenAPI options
		option.Summary("Echo controller with options"),
		option.Description("Some description"),
		option.OperationID("MyCustomOperationID"),
		option.Tags("Echo"),
	)

	// 3️⃣ Level 3: Register Fuego controller to Echo router. Fuego takes care of serialization/deserialization, error handling, content negotiation, etc.
	fuegoecho.Get(engine, echoRouter, "/fuego", fuegoControllerGet)

	// 4️⃣ Level 4: Add some options to the POST endpoint (checks at start-time + validations at request time)
	fuegoecho.Post(engine, echoRouter, "/fuego-with-options", fuegoControllerPost,
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

	// TODO: Supports groups & path parameters even for Echo handlers
	// group := echoRouter.Group("/my-group/:id")
	// fuegoecho.Get(engine, group, "/fuego", fuegoControllerGet,
	// 	option.Summary("Route with group and id"),
	// 	option.Tags("Fuego"),
	// )

	// Serve the OpenAPI spec

	engine.RegisterOpenAPIRoutes(&fuegoecho.OpenAPIHandler{Echo: echoRouter})

	return echoRouter, engine.OpenAPI
}

func (h *HelloRequest) InTransform(ctx context.Context) error {
	// Transformation
	h.Word = strings.ToLower(h.Word)

	// Custom validation, with fuego provided error
	if h.Word == "apple" {
		return fuego.BadRequestError{Title: "Word not allowed", Err: errors.New("forbidden word"), Detail: "The word 'apple' is not allowed"}
	}

	// Custom validation, with basic error
	if h.Word == "banana" {
		return errors.New("banana is not allowed")
	}

	// Context-based transformation
	if user := ctx.Value("user"); user == "secret agent" {
		h.Word = "*****"
	}

	return nil
}
