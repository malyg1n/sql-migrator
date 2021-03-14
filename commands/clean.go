package commands

import (
	"database/sql"
	"github.com/malyg1n/sql-migrator/output"
	"strings"
)

type CleanCommand struct {
	AbstractCommand
}

func NewCleanCommand(db *sql.DB) *CleanCommand {
	return &CleanCommand{
		AbstractCommand{
			db: db,
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
	migrations, err := c.getMigrationsFromBD()
	if err != nil {
		output.ShowError(err.Error())
		return exitStatusError
	}

	for _, me := range migrations {
		if err := c.rollbackMigration(me); err != nil {
			output.ShowError(err.Error())
			return exitStatusError
		}
	}

	return exitStatusSuccess
}
