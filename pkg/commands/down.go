package commands

import (
	"fmt"
	"github.com/malyg1n/sql-migrator/pkg/output"
	"github.com/malyg1n/sql-migrator/pkg/services"
	"strings"
)

type DownCommand struct {
	AbstractCommand
}

func NewDownCommand(service services.ServiceContract) *DownCommand {
	return &DownCommand{
		AbstractCommand{
			service: service,
		},
	}
}

func (c *DownCommand) Help() string {
	helpText := `
Usage: sql-migrator down [directory]
  Undo a database migration.
Options:
  directory			     Directory with migration files (default migrations).
`
	return strings.TrimSpace(helpText)
}

func (c *DownCommand) Synopsis() string {
	return "Undo a database migration."
}

func (c *DownCommand) Run(args []string) int {
	rolledBack, err := c.service.ApplyMigrationsDown()
	if err != nil {
		output.ShowError(err.Error())
		return exitStatusError
	}

	for _, rb := range rolledBack {
		output.ShowWarning(fmt.Sprintf("rolled back: %s", rb))
	}

	return exitStatusSuccess
}
