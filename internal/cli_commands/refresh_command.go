package cli_commands

import (
	"github.com/malyg1n/sql-migrator/internal"
	"strings"
)

type RefreshCommand struct {
	AbstractCommand
}

func NewRefreshCommand(service internal.ServiceContract) *RefreshCommand {
	return &RefreshCommand{
		AbstractCommand{
			service: service,
		},
	}
}

func (c *RefreshCommand) Help() string {
	helpText := `
Usage: sql-migrator refresh [directory]
  Refresh all migrations.
Options:
  directory              Directory with migration files (default migrations).
`
	return strings.TrimSpace(helpText)
}

func (c *RefreshCommand) Synopsis() string {
	return "Refresh all migrations."
}

func (c *RefreshCommand) Run(args []string) int {
	messages, err := c.service.RefreshMigrations()
	if err != nil {
		ShowError(err.Error())
		return exitStatusError
	}

	for _, m := range messages {
		ShowMessage(m)
	}

	return exitStatusSuccess
}
