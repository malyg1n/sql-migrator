package cli_commands

import (
	"flag"
	"fmt"
	"github.com/malyg1n/sql-migrator/internal"
	"strings"
)

// Command for create migrations files
type CreateCommand struct {
	AbstractCommand
}

var (
	migrationName string
)

// Return command instance
func NewCreateCommand(service internal.ServiceContract) *CreateCommand {
	return &CreateCommand{
		AbstractCommand{
			service: service,
		},
	}
}

// Show help text
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

// Show info about command
func (c *CreateCommand) Synopsis() string {
	return "Create a new migration."
}

// Execute command
func (c *CreateCommand) Run(args []string) int {
	flags := flag.NewFlagSet("create", flag.ContinueOnError)
	flags.Parse(args)
	migrationName = flags.Arg(0)

	if migrationName == "" {
		internal.ShowError("empty migration name")
		return exitStatusError
	}

	files, err := c.service.CreateMigrationFile(migrationName)

	if err != nil {
		internal.ShowError(err.Error())
		return exitStatusError
	}

	for _, f := range files {
		internal.ShowMessage(fmt.Sprintf("created empty migration in: %s", f))
	}

	return exitStatusSuccess
}
