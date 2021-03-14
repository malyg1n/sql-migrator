package commands

type CreateCommand struct {
}

func (c *CreateCommand) Help() string {
	return ""
}

func (c *CreateCommand) Synopsis() string {
	return ""
}

func (c *CreateCommand) Run(args []string) int {
	return 0
}
