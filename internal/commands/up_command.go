package commands

import (
	"fmt"
	"github.com/malyg1n/sql-migrator/internal/output"
	"strings"
)

// UpCommand is command for roll out migrations
type UpCommand struct {
	AbstractCommand
}

// NewUpCommand returns command instance
func NewUpCommand(service serviceContract) *UpCommand {
	return &UpCommand{
		AbstractCommand{
			service: service,
		},
	}
}

// Help method displays info about command
func (c *UpCommand) Help() string {
	helpText := `
Usage: sql-migrator up [directory]
  Migrates the database to the most recent version available.
Options:
  directory			     Directory with migration files (default migrations).
`
	return strings.TrimSpace(helpText)
}

// Synopsis method show short description about command
func (c *UpCommand) Synopsis() string {
	return "Migrates the database to the most recent version available."
}

// Run method executes the command
func (c *UpCommand) Run(args []string) int {
	migrated, err := c.service.ApplyMigrationsUp()
	console := output.NewConsoleOutput()
	if err != nil {
		console.PrintError(err.Error())
		return exitStatusError
	}

	if migrated == nil {
		console.PrintInfo("nothing to migrate")
	} else {
		for _, m := range migrated {
			console.PrintSuccess(fmt.Sprintf("migrated: %s", m))
		}
	}

	return exitStatusSuccess
}
