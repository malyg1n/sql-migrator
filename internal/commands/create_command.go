package commands

import (
	"flag"
	"fmt"
	"github.com/malyg1n/sql-migrator/internal/output"
	"strings"
)

// CreateCommand is command for create migrations files
type CreateCommand struct {
	AbstractCommand
}

var (
	migrationName string
)

// NewCreateCommand returns command instance
func NewCreateCommand(service serviceContract) *CreateCommand {
	return &CreateCommand{
		AbstractCommand{
			service: service,
		},
	}
}

// Help method displays info about command
func (c *CreateCommand) Help() string {
	helpText := `
Usage: sql-migrator create [directory] name
  Create a new a database migration.
Options:
  directory              The name of the migrations' directory
  name                   The name of the migration
`
	return strings.TrimSpace(helpText)
}

// Synopsis method show short description about command
func (c *CreateCommand) Synopsis() string {
	return "Create a new migration."
}

// Run method executes the command
func (c *CreateCommand) Run(args []string) int {
	flags := flag.NewFlagSet("create", flag.ContinueOnError)
	err := flags.Parse(args)
	console := output.NewConsoleOutput()
	if err != nil {
		console.PrintError(err.Error())
		return exitStatusError
	}

	migrationName = flags.Arg(0)
	if migrationName == "" {
		console.PrintError("empty migration name")
		return exitStatusError
	}

	files, err := c.service.CreateMigrationFiles(migrationName)

	if err != nil {
		console.PrintError(err.Error())
		return exitStatusError
	}

	for _, f := range files {
		console.PrintSuccess(fmt.Sprintf("created empty migration in: %s", f))
	}

	return exitStatusSuccess
}
