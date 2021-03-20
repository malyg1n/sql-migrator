package commands

import (
	"github.com/malyg1n/sql-migrator/internal/output"
	"strings"
)

// RefreshCommand is command for cleaning of all migrations and roll out them over again
type RefreshCommand struct {
	AbstractCommand
}

// NewRefreshCommand returns command instance
func NewRefreshCommand(service serviceContract) *RefreshCommand {
	return &RefreshCommand{
		AbstractCommand{
			service: service,
		},
	}
}

// Help method displays info about command
func (c *RefreshCommand) Help() string {
	helpText := `
Usage: sql-migrator refresh [directory]
  Refresh all migrations.
Options:
  directory              Directory with migration files (default migrations).
`
	return strings.TrimSpace(helpText)
}

// Synopsis method show short description about command
func (c *RefreshCommand) Synopsis() string {
	return "Refresh all migrations."
}

// Run method executes the command
func (c *RefreshCommand) Run(args []string) int {
	messages, err := c.service.RefreshMigrations()
	console := output.NewConsoleOutput()
	if err != nil {
		console.PrintError(err.Error())
		return exitStatusError
	}

	for _, m := range messages {
		console.PrintSuccess(m)
	}

	return exitStatusSuccess
}
