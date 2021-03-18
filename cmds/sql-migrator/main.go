package main

import (
	"database/sql"
	"github.com/joho/godotenv"
	"github.com/malyg1n/sql-migrator/internal/cli_commands"
	"github.com/malyg1n/sql-migrator/internal/output"
	"github.com/malyg1n/sql-migrator/internal/service"
	"github.com/malyg1n/sql-migrator/internal/store"
	"os"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
	"github.com/mitchellh/cli"
)

const migrationsTableName = "schema_migrations"

// main method
func main() {
	err := godotenv.Load()
	console := output.NewConsoleOutput()
	if err != nil {
		console.PrintError("error loading environment file")
		os.Exit(1)
	}

	db, err := sql.Open(os.Getenv("DB_DRIVER"), os.Getenv("DB_DSN"))
	defer db.Close()

	if err != nil {
		console.PrintError(err.Error())
		os.Exit(1)
	}

	store := store.NewStore(db, "schema_migrations", os.Getenv("DB_DRIVER"))
	service := service.NewService(store, os.Getenv("MIGRATIONS_PATH"))

	c := cli.NewCLI("migrator", "0.0.5")
	c.Args = os.Args[1:]
	c.Commands = map[string]cli.CommandFactory{
		"init": func() (cli.Command, error) {
			return cli_commands.NewInitCommand(service), nil
		},
		"create": func() (cli.Command, error) {
			return cli_commands.NewCreateCommand(service), nil
		},
		"up": func() (cli.Command, error) {
			return cli_commands.NewUpCommand(service), nil
		},
		"down": func() (cli.Command, error) {
			return cli_commands.NewDownCommand(service), nil
		},
		"refresh": func() (cli.Command, error) {
			return cli_commands.NewRefreshCommand(service), nil
		},
		"clean": func() (cli.Command, error) {
			return cli_commands.NewCleanCommand(service), nil
		},
	}

	status, err := c.Run()
	if err != nil {
		console.PrintError(err.Error())
		os.Exit(1)
	}

	os.Exit(status)
}
