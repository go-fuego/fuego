package commands

import (
	"fmt"
	"os"
	"strings"

	"github.com/go-fuego/fuego/cmd/fuego/templates"
	"github.com/urfave/cli/v2"
)

func ControllerCommand() *cli.Command {
	return &cli.Command{
		Name:    "controller",
		Usage:   "add a new template",
		Aliases: []string{"c"},
		Action: func(cCtx *cli.Context) error {
			controllerName := "newController"
			if cCtx.NArg() > 0 {
				controllerName = cCtx.Args().First()
			} else {
				fmt.Println("Note: You can add a controller name as an argument. Example: fuego controller yourControllerName")
			}

			err := createController(controllerName)
			if err != nil {
				return err
			}

			fmt.Printf("🔥 Controller %s created successfully\n", controllerName)
			return nil
		},
	}
}

func createController(controllerName string) error {
	controllerDir := "./controllers/"
	if _, err := os.Stat(controllerDir); os.IsNotExist(err) {
		err = os.Mkdir(controllerDir, 0755)
		if err != nil {
			return err
		}
	}

	templateContent, err := templates.FS.ReadFile("controller.template")

	newContent := strings.ReplaceAll(string(templateContent), "newController", controllerName)
	newContent = strings.ReplaceAll(newContent, "NewController", strings.Title(controllerName))

	controllerPath := fmt.Sprintf("%s%s.go", controllerDir, controllerName)
	err = os.WriteFile(controllerPath, []byte(newContent), 0644)
	if err != nil {
		return err
	}

	return nil
}
