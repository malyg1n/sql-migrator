package commands

import (
	"errors"
	"flag"
	"fmt"
	"github.com/malyg1n/sqlx-migrator/output"
	"os"
	"path"
	"strings"
	"time"
)

// Command for create migrations files
type CreateCommand struct {
}

var (
	migrationName string
)

// Return command instance
func NewCreateCommand() *CreateCommand {
	return &CreateCommand{}
}

// Show help text
func (c *CreateCommand) Help() string {
	helpText := `
Usage: sqlx-migrator create [directory] name
  Create a new a database migration.
Options:
  directory              The name of the migrations' directory
  name                   The name of the migration
`
	return strings.TrimSpace(helpText)
}

// Show info about command
func (c *CreateCommand) Synopsis() string {
	return "Create a new migration"
}

// Execute command
func (c *CreateCommand) Run(args []string) int {
	flags := flag.NewFlagSet("create", flag.ContinueOnError)
	flags.Parse(args)

	if err := c.parseFlags(args); err != nil {
		output.ShowError(err.Error())
		return exitStatusError
	}

	if err := c.createMigration(); err != nil {
		output.ShowError(err.Error())
		return exitStatusError
	}

	return exitStatusSuccess
}

func (c *CreateCommand) createMigration() error {
	upFileName := fmt.Sprintf("%s-%s-up.sql", time.Now().Format(timeFormat), strings.TrimSpace(migrationName))
	pathName := path.Join(migrationDir, upFileName)
	fUp, err := os.Create(pathName)

	if err != nil {
		return err
	}

	output.ShowMessage(fmt.Sprintf("created migration %s", pathName))

	downFileName := fmt.Sprintf("%s-%s-down.sql", time.Now().Format(timeFormat), strings.TrimSpace(migrationName))
	pathName = path.Join(migrationDir, downFileName)
	fDown, err := os.Create(pathName)

	if err != nil {
		return err
	}

	output.ShowMessage(fmt.Sprintf("created migration %s", pathName))

	defer func() {
		_ = fUp.Close()
		_ = fDown.Close()
	}()

	return nil
}

func (c *CreateCommand) parseFlags(args []string) error {
	flags := flag.NewFlagSet("create", flag.ContinueOnError)
	flags.Parse(args)

	switch len(args) {
	case 2:
		migrationDir = flags.Arg(0)
		migrationName = flags.Arg(1)
	case 1:
		migrationName = flags.Arg(1)
	default:
		return errors.New("migration name is required")
	}
	return nil
}
