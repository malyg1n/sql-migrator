package commands

import (
	"fmt"
	"github.com/malyg1n/sql-migrator/pkg/output"
	"github.com/malyg1n/sql-migrator/pkg/services"
	"strings"
)

type CleanCommand struct {
	AbstractCommand
}

func NewCleanCommand(service services.ServiceContract) *CleanCommand {
	return &CleanCommand{
		AbstractCommand{
			service: service,
		},
	}
}

func (c *CleanCommand) Help() string {
	helpText := `
Usage: sql-migrator clean
  Down all migrations.
`
	return strings.TrimSpace(helpText)
}

func (c *CleanCommand) Synopsis() string {
	return "Down all migrations."
}

func (c *CleanCommand) Run(args []string) int {
	rollerBack, err := c.service.ApplyAllMigrationsDown()
	if err != nil {
		output.ShowError(err.Error())
		return exitStatusError
	}

	for _, rb := range rollerBack {
		output.ShowWarning(fmt.Sprintf("rolled back: %s", rb))
	}

	return exitStatusSuccess
}
