package commands

import (
	"database/sql"
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
	return ""
}

func (c *DownCommand) Synopsis() string {
	return ""
}

func (c *DownCommand) Run(args []string) int {
	return 0
}
