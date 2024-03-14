package openapi3testing

import (
	"testing"

	"github.com/go-fuego/fuego"
)

func TestValidateBasicDocument(t *testing.T) {
	s := fuego.NewServer()

	fuego.Get(s, "/", func(c *fuego.ContextNoBody) (string, error) {
		return "Hello, World!", nil
	})

	validateSpec(s)
}

type MyInput struct {
	Name string `json:"name" validate:"required,max=30"`
}

type MyOutput struct {
	Answer string `json:"answer"`
}

func TestValidateDocumentWithParams(t *testing.T) {
	s := fuego.NewServer()

	fuego.Get(s, "/:id", func(c *fuego.ContextWithBody[MyInput]) (MyOutput, error) {
		return MyOutput{Answer: "Hello " + c.PathParam("id")}, nil
	})

	validateSpec(s)
}
