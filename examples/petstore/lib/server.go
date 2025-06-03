package lib

import (
	"net/http"
	"reflect"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/getkin/kin-openapi/openapi3gen"
	"github.com/go-fuego/fuego"
	controller "github.com/go-fuego/fuego/examples/petstore/controllers"
	"github.com/go-fuego/fuego/examples/petstore/services"
	"github.com/go-fuego/fuego/option"
)

type NoContent struct {
	Empty string `json:"-"`
}

var uuidParser openapi3gen.SchemaCustomizerFn = func(name string, t reflect.Type, tag reflect.StructTag, schema *openapi3.Schema) error {
	validateTag := tag.Get("validate")
	if validateTag == "" {
		return nil
	}

	parts := strings.Split(validateTag, ",")
	for _, part := range parts {
		if part == "uuid" {
			schema.Format = "uuid"
		}
	}
	return nil
}

func NewPetStoreServer(options ...func(*fuego.Server)) *fuego.Server {
	options = append(options, fuego.WithRouteOptions(
		option.AddResponse(http.StatusNoContent, "No Content", fuego.Response{Type: NoContent{}}),
	))
	options = append(options, fuego.WithEngineOptions(
		fuego.WithOpenAPIGeneratorSchemaCustomizer(
			uuidParser,
		),
	))

	s := fuego.NewServer(options...)

	petsResources := controller.PetsResources{
		PetsService: services.NewInMemoryPetsService(), // Dependency injection: we can pass a service here (for example a database service)
	}
	petsResources.Routes(s)
	return s
}
