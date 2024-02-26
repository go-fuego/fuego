package main

import (
	"net/http"

	"github.com/go-fuego/fuego"
	jsoniter "github.com/json-iterator/go"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

func main() {
	s := fuego.NewServer()

	s.Serialize = func(w http.ResponseWriter, ans any) {
		w.Header().Set("Content-Type", "text/plain")
		json.NewEncoder(w).Encode(ans)
	}

	fuego.Get(s, "/", helloWorld)

	s.Run()
}

func helloWorld(c fuego.ContextNoBody) (string, error) {
	return "Hello, World!", nil
}
