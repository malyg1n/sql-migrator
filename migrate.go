package main

import (
	"fmt"
	"github.com/malyg1n/sqlx-migrator.git/commands"
	"os"

	"github.com/mitchellh/cli"
)

func main() {
	c := cli.NewCLI("migrate", "0.0.1")
	c.Args = os.Args[1:]
	c.Commands = map[string]cli.CommandFactory{
		"create": func() (cli.Command, error) {
			return &commands.CreateCommand{}, nil
		},
		"up": func() (cli.Command, error) {
			return &commands.CreateCommand{}, nil
		},
		"down": func() (cli.Command, error) {
			return &commands.CreateCommand{}, nil
		},
	}
	status, err := c.Run()
	if err != nil {
		fmt.Println(err)
	}
	os.Exit(status)
}
