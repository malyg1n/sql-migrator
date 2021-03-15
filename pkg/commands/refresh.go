package commands

import (
	"github.com/malyg1n/sql-migrator/pkg/output"
	"github.com/malyg1n/sql-migrator/pkg/services"
	"strings"
)

type RefreshCommand struct {
	AbstractCommand
}

func NewRefreshCommand(service services.ServiceContract) *RefreshCommand {
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
		output.ShowError(err.Error())
		return exitStatusError
	}

	for _, m := range messages {
		output.ShowMessage(m)
	}

	return exitStatusSuccess
}
