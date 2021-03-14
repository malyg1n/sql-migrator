package commands

import (
	"database/sql"
	"flag"
	"github.com/malyg1n/sql-migrator/output"
	"strings"
)

type UpCommand struct {
	AbstractCommand
}

func NewUpCommand(db *sql.DB) *UpCommand {
	return &UpCommand{
		AbstractCommand{
			db: db,
		},
	}
}

func (c *UpCommand) Help() string {
	helpText := `
Usage: sql-migrator up [directory]
  Migrates the database to the most recent version available.
Options:
  directory			     Directory with migration files (default migrations).
`
	return strings.TrimSpace(helpText)
}

func (c *UpCommand) Synopsis() string {
	return "Migrates the database to the most recent version available."
}

func (c *UpCommand) Run(args []string) int {
	flags := flag.NewFlagSet("up", flag.ContinueOnError)
	flags.Parse(args)
	if mDir := flags.Arg(0); mDir != "" {
		migrationDir = mDir
	}
	if err := c.checkFolder(migrationDir); err != nil {
		output.ShowError(err.Error())
		return exitStatusError
	}

	if err := c.checkMigrationsTable(); err != nil {
		output.ShowError(err.Error())
		return exitStatusError
	}

	migrations, err := c.getMigrationsFromBD()
	if err != nil {
		output.ShowError(err.Error())
		return exitStatusError
	}

	files, err := c.getMigrationUpFiles(migrationDir)
	if err != nil {
		output.ShowError(err.Error())
		return exitStatusError
	}

	newMigrationsFiles := c.filterMigrations(migrations, files)
	if len(newMigrationsFiles) > 0 {
		vn, err := c.getLatestVersionNumber()
		if err != nil {
			output.ShowError(err.Error())
			return exitStatusError
		}
		vn++
		if err := c.applyMigrationsUp(newMigrationsFiles, vn); err != nil {
			output.ShowError(err.Error())
			return exitStatusError
		}
	} else {
		output.ShowInfo("nothing to migrate")
	}

	return exitStatusSuccess
}
