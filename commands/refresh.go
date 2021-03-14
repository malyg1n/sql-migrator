package commands

import (
	"database/sql"
	"flag"
	"github.com/malyg1n/sql-migrator/output"
	"strings"
)

type RefreshCommand struct {
	AbstractCommand
}

func NewRefreshCommand(db *sql.DB) *RefreshCommand {
	return &RefreshCommand{
		AbstractCommand{
			db: db,
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
	flags := flag.NewFlagSet("up", flag.ContinueOnError)
	flags.Parse(args)
	if mDir := flags.Arg(0); mDir != "" {
		migrationDir = mDir
	}
	if err := c.checkFolder(migrationDir); err != nil {
		output.ShowError(err.Error())
		return exitStatusError
	}
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

	migrationFiles, err := c.getMigrationUpFiles(migrationDir)
	if err != nil {
		output.ShowError(err.Error())
		return exitStatusError
	}

	err = c.applyMigrationsUp(migrationFiles, 1)
	if err != nil {
		output.ShowError(err.Error())
		return exitStatusError
	}

	return exitStatusSuccess
}
