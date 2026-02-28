package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/gorilla/mux"

	"github.com/go-fuego/fuego"
	"github.com/go-fuego/fuego/extra/fuegomux"
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
	r, _ := server()

	fmt.Println("OpenAPI at http://localhost:8980/swagger")

	err := http.ListenAndServe(":8980", r)
	if err != nil {
		panic(err)
	}
}

func server() (*mux.Router, *fuego.OpenAPI) {
	muxRouter := mux.NewRouter()
	engine := fuego.NewEngine()

	// Register native mux controller
	muxRouter.HandleFunc("/mux", muxController).Methods(http.MethodGet)

	// 1. Level 1: Register native mux controller with OpenAPI spec
	fuegomux.GetMux(engine, muxRouter, "/mux-with-openapi", muxController)

	// 2. Level 2: Native mux controller with OpenAPI options
	fuegomux.GetMux(engine, muxRouter, "/mux-with-openapi-and-options", muxController,
		option.Summary("Mux controller with options"),
		option.Description("Some description"),
		option.OperationID("MyCustomOperationID"),
		option.Tags("Mux"),
	)

	// 3. Level 3: Fuego controller with gorilla/mux router
	fuegomux.Get(engine, muxRouter, "/fuego", fuegoControllerGet)

	// 4. Level 4: Fuego controller with options
	fuegomux.Post(engine, muxRouter, "/fuego-with-options", fuegoControllerPost,
		option.Description("Some description"),
		option.OperationID("SomeOperationID"),
		option.AddError(409, "Name Already Exists"),
		option.DefaultStatusCode(201),
		option.Tags("Fuego"),
		option.Query("name", "Your name", param.Example("name example", "John Carmack")),
		option.Header("X-Request-ID", "Request ID", param.Default("123456")),
		option.Header("Content-Type", "Content Type", param.Default("application/json")),
	)

	// Groups & path parameters with regex
	sub := muxRouter.PathPrefix("/my-group/{id:[0-9]+}").Subrouter()
	fuegomux.Get(engine, sub, "/fuego", fuegoControllerGet,
		option.Summary("Route with subrouter and id"),
		option.Tags("Fuego"),
	)

	engine.RegisterOpenAPIRoutes(&fuegomux.OpenAPIHandler{Router: muxRouter})

	return muxRouter, engine.OpenAPI
}

func (h *HelloRequest) InTransform(ctx context.Context) error {
	h.Word = strings.ToLower(h.Word)

	if h.Word == "apple" {
		return fuego.BadRequestError{Title: "Word not allowed", Err: errors.New("forbidden word"), Detail: "The word 'apple' is not allowed"}
	}

	if h.Word == "banana" {
		return errors.New("banana is not allowed")
	}

	if user := ctx.Value("user"); user == "secret agent" {
		h.Word = "*****"
	}

	return nil
}
