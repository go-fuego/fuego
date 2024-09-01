package main

import (
	"net/http"

	jsoniter "github.com/json-iterator/go"

	"github.com/go-fuego/fuego"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

func main() {
	s := fuego.NewServer()

	s.Serialize = func(w http.ResponseWriter, _ *http.Request, ans any) error {
		w.Header().Set("Content-Type", "text/plain")
		return json.NewEncoder(w).Encode(ans)
	}

	fuego.Get(s, "/", helloWorld)

	s.Run()
}

func helloWorld(c fuego.ContextNoBody) (string, error) {
	return "Hello, World!", nil
}
