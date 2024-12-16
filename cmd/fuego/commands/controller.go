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
			entityName := cCtx.Args().First()

			if entityName == "" {
				entityName = "newEntity"
				fmt.Println("Note: You can add a controller name as an argument. Example: `fuego controller books`")
			}

			_, err := createController(entityName)
			if err != nil {
				return err
			}

			fmt.Printf("🔥 Controller %s created successfully\n", entityName)
			return nil
		},
	}
}

// createController creates a new controller file
func createController(entityName string) (string, error) {
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

	newContent := strings.ReplaceAll(string(templateContent), "newEntity", entityName)
	newContent = strings.ReplaceAll(newContent, "NewEntity", titler.String(entityName))

	controllerPath := fmt.Sprintf("%s%s.go", controllerDir, entityName)

	err = os.WriteFile(controllerPath, []byte(newContent), 0o644)
	if err != nil {
		return "", err
	}

	return newContent, nil
}
