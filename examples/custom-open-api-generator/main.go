package main

import (
	"reflect"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/getkin/kin-openapi/openapi3gen"

	"github.com/go-fuego/fuego"
)

type Input struct {
	Name string `json:"name" description:"some name" readOnly:"true"`
}

type Output struct {
	Input Input `json:"input"`
}

func main() {
	s := fuego.NewServer(
		fuego.WithOpenAPIGenerator(
			openapi3gen.NewGenerator(
				openapi3gen.UseAllExportedFields(),
				openapi3gen.SchemaCustomizer(func(name string, t reflect.Type, tag reflect.StructTag, schema *openapi3.Schema) error {
					if v := tag.Get("readOnly"); v == "true" {
						schema.ReadOnly = true
					}

					if v := tag.Get("description"); v != "" {
						schema.Description = v
					}

					return nil
				}),
			),
		),
	)

	fuego.Post(s, "/", func(c *fuego.ContextWithBody[Input]) (Output, error) {
		input, err := c.Body()
		if err != nil {
			return Output{}, err
		}

		return Output{
			Input: input,
		}, nil
	})

	s.Run()
}
