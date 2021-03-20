package commands

import (
	"github.com/malyg1n/sql-migrator/internal/output"
	"strings"
)

// DownCommand is command for rollback last migrations
type DownCommand struct {
	AbstractCommand
}

// NewDownCommand returns command instance
func NewDownCommand(service serviceContract) *DownCommand {
	return &DownCommand{
		AbstractCommand{
			service: service,
		},
	}
}

// Help method displays info about command
func (c *DownCommand) Help() string {
	helpText := `
Usage: sql-migrator down [directory]
  Undo a database migration.
Options:
  directory			     Directory with migration files (default migrations).
`
	return strings.TrimSpace(helpText)
}

// Synopsis method show short description about command
func (c *DownCommand) Synopsis() string {
	return "Undo a database migration."
}

// Run method executes the command
func (c *DownCommand) Run(args []string) int {
	rolledBack, err := c.service.ApplyMigrationsDown()
	console := output.NewConsoleOutput()
	if err != nil {
		console.PrintError(err.Error())
		return exitStatusError
	}

	for _, rb := range rolledBack {
		console.PrintWarning(rb)
	}

	return exitStatusSuccess
}
