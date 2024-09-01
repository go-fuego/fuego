package commands

import (
	"fmt"
	"os"
	"strings"

	"github.com/urfave/cli/v2"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"

	"github.com/go-fuego/fuego/cmd/fuego/templates"
)

func Controller() *cli.Command {
	return &cli.Command{
		Name:    "controller",
		Usage:   "creates a new controller file",
		Aliases: []string{"c"},
		Action: func(cCtx *cli.Context) error {
			controllerName := cCtx.Args().First()

			if controllerName == "" {
				controllerName = "newController"
				fmt.Println("Note: You can add a controller name as an argument. Example: `fuego controller books`")
			}

			_, err := createController(controllerName)
			if err != nil {
				return err
			}

			fmt.Printf("🔥 Controller %s created successfully\n", controllerName)
			return nil
		},
	}
}

// createController creates a new controller file
func createController(controllerName string) (string, error) {
	controllerDir := "./controller/"
	if _, err := os.Stat(controllerDir); os.IsNotExist(err) {
		err = os.Mkdir(controllerDir, 0o755)
		if err != nil {
			return "", err
		}
	}

	templateContent, err := templates.FS.ReadFile("controller/controller.go")
	if err != nil {
		return "", err
	}

	t := language.English
	titler := cases.Title(t)

	newContent := strings.ReplaceAll(string(templateContent), "newController", controllerName)
	newContent = strings.ReplaceAll(newContent, "NewController", titler.String(controllerName))

	controllerPath := fmt.Sprintf("%s%s.go", controllerDir, controllerName)

	err = os.WriteFile(controllerPath, []byte(newContent), 0o644)
	if err != nil {
		return "", err
	}

	return newContent, nil
}
