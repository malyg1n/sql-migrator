package commands

import (
	"fmt"
	"github.com/malyg1n/sql-migrator/internal/output"
	"strings"
)

// CleanCommand is command for clean database
type CleanCommand struct {
	AbstractCommand
}

// NewCleanCommand returns command instance
func NewCleanCommand(service serviceContract) *CleanCommand {
	return &CleanCommand{
		AbstractCommand{
			service: service,
		},
	}
}

// Help method displays info about command
func (c *CleanCommand) Help() string {
	helpText := `
Usage: sql-migrator clean
  Down all migrations.
`
	return strings.TrimSpace(helpText)
}

// Synopsis method show short description about command
func (c *CleanCommand) Synopsis() string {
	return "Down all migrations."
}

// Run method executes the command
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
