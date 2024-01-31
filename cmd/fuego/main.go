package main

import (
	"fmt"
	"log"
	"os"

	"github.com/go-fuego/fuego/cmd/fuego/commands"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:  "fuego",
		Usage: "Fire like fuego!",
		Action: func(*cli.Context) error {
			fmt.Println("The 🔥 CLI!")
			return nil
		},
		Commands: []*cli.Command{
			commands.ControllerCommand(),
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
