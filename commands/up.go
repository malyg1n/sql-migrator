package commands

import "github.com/malyg1n/sqlx-migrator/configs"

type UpCommand struct {
	AbstractCommand
}

func NewUpCommand(cfg *configs.DBConfig) *UpCommand {
	return &UpCommand{
		AbstractCommand{
			config: cfg,
		},
	}
}

func (c *UpCommand) Help() string {
	return ""
}

func (c *UpCommand) Synopsis() string {
	return ""
}

func (c *UpCommand) Run(args []string) int {
	return 0
}
