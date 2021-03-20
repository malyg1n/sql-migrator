package cli_commands

import (
	"fmt"
	"github.com/malyg1n/sql-migrator/internal/output"
	"strings"
)

type CleanCommand struct {
	AbstractCommand
}

func NewCleanCommand(service serviceContract) *CleanCommand {
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
	console := output.NewConsoleOutput()
	if err != nil {
		console.PrintError(err.Error())
		return exitStatusError
	}

	for _, rb := range rollerBack {
		console.PrintWarning(fmt.Sprintf("rolled back: %s", rb))
	}

	return exitStatusSuccess
}
