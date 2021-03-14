package commands

import (
	"database/sql"
	"flag"
	"github.com/malyg1n/sql-migrator/output"
	"strings"
)

type DownCommand struct {
	AbstractCommand
}

func NewDownCommand(db *sql.DB) *DownCommand {
	return &DownCommand{
		AbstractCommand{
			db: db,
		},
	}
}

func (c *DownCommand) Help() string {
	helpText := `
Usage: sql-migrator down [directory]
  Undo a database migration.
Options:
  directory			     Directory with migration files (default migrations).
`
	return strings.TrimSpace(helpText)
}

func (c *DownCommand) Synopsis() string {
	return "Undo a database migration."
}

func (c *DownCommand) Run(args []string) int {
	flags := flag.NewFlagSet("up", flag.ContinueOnError)
	flags.Parse(args)
	if mDir := flags.Arg(0); mDir != "" {
		migrationDir = mDir
	}
	if err := c.checkFolder(migrationDir); err != nil {
		output.ShowError(err.Error())
		return exitStatusError
	}

	version, err := c.getLatestVersionNumber()
	if err != nil {
		output.ShowError(err.Error())
		return exitStatusError
	}

	migrations, err := c.getMigrationsByVersion(version)

	for _, me := range migrations {
		if err := c.rollbackMigration(me); err != nil {
			output.ShowError(err.Error())
			return exitStatusError
		}
	}

	return exitStatusSuccess
}
