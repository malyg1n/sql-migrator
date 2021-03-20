package main

import (
	"database/sql"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/malyg1n/sql-migrator/internal/commands"
	"github.com/malyg1n/sql-migrator/internal/config"
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
	handleBreaker()
	loadEnv()

	cfg := config.NewConfig()
	db, err := sql.Open(cfg.DbDriver, cfg.DbConnectionsString)
	defer func() {
		err := db.Close()
		if err != nil {
			fmt.Errorf("db wasn't closed correctly with error: %s", err.Error())
			os.Exit(1)
		}
	}()

	if err != nil {
		fmt.Errorf("db wasn't opened correctly with error: %s", err.Error())
		os.Exit(1)
	}

	store := store.NewStore(db, migrationsTableName)
	service := service.NewService(store, cfg)

	status, err := initCommands(service)
	if err != nil {
		fmt.Errorf("commands wasn't initialized with error: %s", err.Error())
		os.Exit(1)
	}

	os.Exit(status)
}

// Init cli commands
func initCommands(service *service.Service) (int, error) {
	newCLI := cli.NewCLI("sql-migrator", "0.1.1")
	newCLI.Args = os.Args[1:]
	newCLI.Commands = map[string]cli.CommandFactory{
		"init": func() (cli.Command, error) {
			return commands.NewInitCommand(service), nil
		},
		"create": func() (cli.Command, error) {
			return commands.NewCreateCommand(service), nil
		},
		"up": func() (cli.Command, error) {
			return commands.NewUpCommand(service), nil
		},
		"down": func() (cli.Command, error) {
			return commands.NewDownCommand(service), nil
		},
		"refresh": func() (cli.Command, error) {
			return commands.NewRefreshCommand(service), nil
		},
		"clean": func() (cli.Command, error) {
			return commands.NewCleanCommand(service), nil
		},
	}

	return newCLI.Run()
}

// Listen interrupt signal
func handleBreaker() {
	ch := make(chan os.Signal)
	signal.Notify(ch, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-ch
		os.Exit(0)
	}()
}

// Loading environments from .env
func loadEnv() {
	err := godotenv.Load()
	if err != nil {
		fmt.Errorf("error loading environment file with error: %s", err.Error())
		os.Exit(1)
	}
}
