package cli_commands

import (
	"fmt"
	"github.com/malyg1n/sql-migrator/internal"
	"strings"
)

type UpCommand struct {
	AbstractCommand
}

func NewUpCommand(service internal.ServiceContract) *UpCommand {
	return &UpCommand{
		AbstractCommand{
			service: service,
		},
	}
}

func (c *UpCommand) Help() string {
	helpText := `
Usage: sql-migrator up [directory]
  Migrates the database to the most recent version available.
Options:
  directory			     Directory with migration files (default migrations).
`
	return strings.TrimSpace(helpText)
}

func (c *UpCommand) Synopsis() string {
	return "Migrates the database to the most recent version available."
}

func (c *UpCommand) Run(args []string) int {
	migrated, err := c.service.ApplyMigrationsUp()

	if err != nil {
		internal.ShowError(err.Error())
		return exitStatusError
	}

	if migrated == nil {
		internal.ShowInfo("nothing to migrate")
	} else {
		for _, m := range migrated {
			internal.ShowMessage(fmt.Sprintf("migrated: %s", m))
		}
	}

	return exitStatusSuccess
}
