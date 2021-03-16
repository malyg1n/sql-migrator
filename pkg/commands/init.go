package commands

import (
	"github.com/malyg1n/sql-migrator/pkg/services"
	"strings"
)

type InitCommand struct {
	AbstractCommand
}

// Return command instance
func NewInitCommand(service services.ServiceContract) *InitCommand {
	return &InitCommand{
		AbstractCommand{
			service: service,
		},
	}
}

// Show help text
func (c *InitCommand) Help() string {
	helpText := `
Usage: sql-migrator init
  Init a table for store your migrations.
`
	return strings.TrimSpace(helpText)
}

// Show info about command
func (c *InitCommand) Synopsis() string {
	return "Init a table for store your migrations."
}

// Execute command
func (c *InitCommand) Run(args []string) int {
	return exitStatusSuccess
}
