package main

import (
	"errors"
	"log"
	"net/http"
	"strings"

	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/jub0bs/cors"

	"github.com/go-fuego/fuego"
	"github.com/go-fuego/fuego/option"
)

type MyError struct {
	Err     error  `json:"error"`
	Message string `json:"message"`
}

var (
	_ fuego.ErrorWithStatus = MyError{}
	_ fuego.ErrorWithDetail = MyError{}
)

func (e MyError) Error() string { return e.Err.Error() }

func (e MyError) StatusCode() int { return http.StatusTeapot }

func (e MyError) DetailMsg() string {
	return strings.Split(e.Error(), " ")[1]
}

func main() {
	s := fuego.NewServer(
		fuego.WithAddr("localhost:8088"),
	)

	cors, err := cors.NewMiddleware(cors.Config{
		Origins:        []string{"*"},
		Methods:        []string{http.MethodGet, http.MethodHead, http.MethodPost},
		RequestHeaders: []string{"Accept", "Content-Type", "X-Requested-With"},
	})
	if err != nil {
		log.Fatal(err)
	}
	fuego.Use(s, cors.Wrap)

	fuego.Use(s, chiMiddleware.Compress(5, "text/html", "text/css"))

	fuego.Get(s, "/custom-err", func(c fuego.ContextNoBody) (string, error) {
		return "hello", MyError{Err: errors.New("my error")}
	},
		option.AddError(http.StatusTeapot, "my custom teapot error", MyError{}),
	)

	s.Run()
}
