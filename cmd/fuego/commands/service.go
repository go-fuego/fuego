package commands

import (
	"fmt"

	"github.com/urfave/cli/v2"
)

func Service() *cli.Command {
	return &cli.Command{
		Name:    "service",
		Usage:   "creates a new service file",
		Aliases: []string{"s"},
		Action:  serviceCommandAction,
	}
}

func serviceCommandAction(cCtx *cli.Context) error {
	entityName := cCtx.Args().First()

	if entityName == "" {
		entityName = "newController"
		fmt.Println("Note: You can add an entity name as an argument. Example: `fuego service books`")
	}

	_, err := createNewEntityDomainFile(entityName, "entity.go", entityName+".go")
	if err != nil {
		return err
	}

	_, err = createNewEntityDomainFile(entityName, "service.go", entityName+"Service.go")
	if err != nil {
		return err
	}

	fmt.Printf("🔥 Service %s created successfully\n", entityName)
	return nil
}
