package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"github.com/malyg1n/sqlx-migrator/commands"
	"github.com/malyg1n/sqlx-migrator/configs"
	"log"
	"os"

	"github.com/mitchellh/cli"
)

var (
	ui cli.Ui
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading environment file")
	}

	cfg := configs.NewDBConfig()

	c := cli.NewCLI("migrate", "0.0.1")
	c.Args = os.Args[1:]
	c.Commands = map[string]cli.CommandFactory{
		"create": func() (cli.Command, error) {
			return commands.NewCreateCommand(), nil
		},
		"up": func() (cli.Command, error) {
			return commands.NewDownCommand(cfg), nil
		},
		"down": func() (cli.Command, error) {
			return commands.NewDownCommand(cfg), nil
		},
	}
	status, err := c.Run()
	if err != nil {
		fmt.Println(err)
	}
	os.Exit(status)
}
