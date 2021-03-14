package commands

import (
	"database/sql"
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
	return ""
}

func (c *CleanCommand) Synopsis() string {
	return ""
}

func (c *CleanCommand) Run(args []string) int {
	return 0
}
