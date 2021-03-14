package commands

import "github.com/malyg1n/sqlx-migrator/configs"

type DownCommand struct {
	AbstractCommand
}

func NewDownCommand(cfg *configs.DBConfig) *DownCommand {
	return &DownCommand{
		AbstractCommand{
			config: cfg,
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
