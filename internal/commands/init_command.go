package commands

import (
	"github.com/malyg1n/sql-migrator/internal/output"
	"strings"
)

// InitCommand is command for initialize utility (create folder and table for migrations)
type InitCommand struct {
	AbstractCommand
}

// NewInitCommand returns command instance
func NewInitCommand(service serviceContract) *InitCommand {
	return &InitCommand{
		AbstractCommand{
			service: service,
		},
	}
}

// Help method displays info about command
func (c *InitCommand) Help() string {
	helpText := `
Usage: sql-migrator init
  Init a table for store your migrations.
`
	return strings.TrimSpace(helpText)
}

// Synopsis method show short description about command
func (c *InitCommand) Synopsis() string {
	return "Init a table for store your migrations."
}

// Run method executes the command
func (c *InitCommand) Run(args []string) int {
	err := c.service.Prepare()
	console := output.NewConsoleOutput()
	if err != nil {
		console.PrintError(err.Error())
		return exitStatusError
	}

	console.PrintSuccess("migrator was initialized")

	return exitStatusSuccess
}
