package main

import (
	"database/sql"
	"github.com/joho/godotenv"
	"github.com/malyg1n/sql-migrator/internal/cli_commands"
	"github.com/malyg1n/sql-migrator/internal/config"
	"github.com/malyg1n/sql-migrator/internal/output"
	"github.com/malyg1n/sql-migrator/internal/service"
	"github.com/malyg1n/sql-migrator/internal/store"
	"os"
	"os/signal"
	"syscall"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
	"github.com/mitchellh/cli"
)

const migrationsTableName = "schema_migrations"

// main method
func main() {
	ch := make(chan os.Signal)
	signal.Notify(ch, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-ch
		os.Exit(0)
	}()

	err := godotenv.Load()
	console := output.NewConsoleOutput()
	if err != nil {
		console.PrintError("error loading environment file")
		os.Exit(1)
	}

	cfg := config.NewConfig()

	db, err := sql.Open(cfg.DbDriver, cfg.DbConnectionsString)
	defer func() {
		err := db.Close()
		if err != nil {
			console.PrintError("db wasn't closed correctly")
			os.Exit(1)
		}
	}()

	if err != nil {
		console.PrintError(err.Error())
		os.Exit(1)
	}

	store := store.NewStore(db, migrationsTableName)
	service := service.NewService(store, cfg)

	newCLI := cli.NewCLI("migrator", "0.0.5")
	newCLI.Args = os.Args[1:]
	newCLI.Commands = map[string]cli.CommandFactory{
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

	status, err := newCLI.Run()
	if err != nil {
		console.PrintError(err.Error())
		os.Exit(1)
	}

	os.Exit(status)
}
