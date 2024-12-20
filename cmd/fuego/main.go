package main

import (
	"fmt"
	"log"
	"os"

	"github.com/urfave/cli/v2"

	"github.com/go-fuego/fuego/cmd/fuego/commands"
)

func main() {
	app := &cli.App{
		Name:  "fuego",
		Usage: "The framework for busy Go developers",
		Action: func(c *cli.Context) error {
			fmt.Println("The 🔥 CLI!")
			return nil
		},
		Commands: []*cli.Command{
			commands.Controller(),
			commands.Service(),
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
