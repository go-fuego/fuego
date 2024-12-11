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

			_, err := createControllerFile(controllerName, "controller.go", controllerName+".go")
			if err != nil {
				return err
			}

			fmt.Printf("🔥 Controller %s created successfully\n", controllerName)
			return nil
		},
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:  "with-service",
				Usage: "enable service file generation",
				Value: false,
				Action: func(cCtx *cli.Context, shouldGenerateServiceFile bool) error {
					if !shouldGenerateServiceFile {
						return nil
					}

					controllerName := cCtx.Args().First()

					_, err := createControllerFile(controllerName, "service.go", controllerName+"Service.go")
					if err != nil {
						return err
					}

					return nil
				},
			},
		},
	}
}

// createController creates a new controller file
func createControllerFile(controllerName, controllerTemplateFileName, outputFileName string) (string, error) {
	controllerDir := "./controller/"
	if _, err := os.Stat(controllerDir); os.IsNotExist(err) {
		err = os.Mkdir(controllerDir, 0o755)
		if err != nil {
			return "", err
		}
	}

	templateContent, err := templates.FS.ReadFile("controller/" + controllerTemplateFileName)
	if err != nil {
		return "", err
	}

	t := language.English
	titler := cases.Title(t)

	newContent := strings.ReplaceAll(string(templateContent), "newController", controllerName)
	newContent = strings.ReplaceAll(newContent, "NewController", titler.String(controllerName))

	controllerPath := fmt.Sprintf("%s%s", controllerDir, outputFileName)

	err = os.WriteFile(controllerPath, []byte(newContent), 0o644)
	if err != nil {
		return "", err
	}

	return newContent, nil
}
