package commands

type DownCommand struct {
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
