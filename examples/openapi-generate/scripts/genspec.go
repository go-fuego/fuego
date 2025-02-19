package main

import (
	"encoding/json"
	"github.com/go-fuego/fuego/examples/openapi-generate/server"
	"os"
)

func main() {
	spec := server.GetServer().OutputOpenAPISpec()

	b, err := json.MarshalIndent(spec, "", "  ")
	if err != nil {
		panic(err)
	}

	err = os.MkdirAll("api", 0750)
	if err != nil {
		panic(err)
	}

	err = os.WriteFile("api/openapi.json", b, 0644)
	if err != nil {
		panic(err)
	}
}
