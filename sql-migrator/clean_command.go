package sql_migrator

import (
	"fmt"
	"strings"
)

type CleanCommand struct {
	AbstractCommand
}

func NewCleanCommand(service ServiceContract) *CleanCommand {
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
		ShowError(err.Error())
		return exitStatusError
	}

	for _, rb := range rollerBack {
		ShowWarning(fmt.Sprintf("rolled back: %s", rb))
	}

	return exitStatusSuccess
}
