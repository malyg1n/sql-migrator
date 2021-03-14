package commands

import (
	"database/sql"
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
	return ""
}

func (c *RefreshCommand) Synopsis() string {
	return ""
}

func (c *RefreshCommand) Run(args []string) int {
	return 0
}
