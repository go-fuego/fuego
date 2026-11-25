package main

import (
	"encoding/json"
	"fmt"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/go-fuego/fuego"
)

type Owner struct {
	Name string `json:"name"`
}

type Widget struct {
	Name  *string  `json:"name"`  // pointer -> nullable scalar
	Tags  []string `json:"tags"`  // slice
	Owner *Owner   `json:"owner"` // pointer to a named struct -> $ref
}

func main() {
	s := fuego.NewServer(
		fuego.WithEngineOptions(
			fuego.WithWalkSchemas(func(_ string, ref *openapi3.SchemaRef) error {
				schema := ref.Value
				if schema.Nullable && schema.Type != nil && !schema.Type.Includes("null") {
					*schema.Type = append(*schema.Type, "null")
					schema.Nullable = false
				}
				return nil
			}),
		),
	)

	fuego.Get(s, "/widget", func(c fuego.ContextNoBody) (Widget, error) {
		return Widget{}, nil
	})

	// Owner is also returned on its own, where it is never null.
	fuego.Get(s, "/owner", func(c fuego.ContextNoBody) (Owner, error) {
		return Owner{}, nil
	})

	doc := s.OutputOpenAPISpec()
	fmt.Println("openapi:", doc.OpenAPI)
	for _, name := range []string{"Widget", "Owner"} {
		out, _ := json.MarshalIndent(doc.Components.Schemas[name].Value, "", "  ")
		fmt.Printf("%s: %s\n", name, out)
	}
}
