package server

import (
	"github.com/go-fuego/fuego"
	"github.com/go-fuego/fuego/option"
)

func helloWorld(c fuego.ContextNoBody) (string, error) {
	return "Hello, World!", nil
}

// GetServer is a representation of "central" server configuration
// that can be used in multiple places, e.g. in the main function and in the OpenAPI spec generation script.
func GetServer() *fuego.Server {

	s := fuego.NewServer()
	// Disable local save of the OpenAPI spec after runtime
	s.Engine.OpenAPI.Config.DisableLocalSave = true

	fuego.Get(s, "/", helloWorld,
		option.Summary("A simple hello world"),
		option.Description("This is a simple hello world"),
		option.Deprecated(),
	)

	return s
}
