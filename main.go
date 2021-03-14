package main

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/malyg1n/sql-migrator/commands"
	"github.com/malyg1n/sql-migrator/configs"
	"github.com/malyg1n/sql-migrator/output"
	"os"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
	"github.com/mitchellh/cli"
)

// main method
func main() {
	err := godotenv.Load()
	if err != nil {
		output.ShowError("error loading environment file")
		os.Exit(1)
	}

	cfg := configs.NewDBConfig()
	db, err := initDB(cfg)
	if err != nil {
		output.ShowError(err.Error())
		os.Exit(1)
	}

	status, err := initCommands(db)
	if err != nil {
		output.ShowError(err.Error())
		os.Exit(1)
	}

	os.Exit(status)
}

// Init list of commands
func initCommands(db *sql.DB) (int, error) {
	c := cli.NewCLI("migrate", "0.0.1")
	c.Args = os.Args[1:]
	c.Commands = map[string]cli.CommandFactory{
		"create": func() (cli.Command, error) {
			return commands.NewCreateCommand(), nil
		},
		"up": func() (cli.Command, error) {
			return commands.NewUpCommand(db), nil
		},
		"down": func() (cli.Command, error) {
			return commands.NewDownCommand(db), nil
		},
		"refresh": func() (cli.Command, error) {
			return commands.NewRefreshCommand(db), nil
		},
	}

	return c.Run()
}

// Init connect to database
func initDB(cfg *configs.DBConfig) (*sql.DB, error) {
	dsn, err := getDSN(cfg)
	if err != nil {
		return nil, err
	}

	db, err := sql.Open(cfg.Driver, dsn)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

// Get dsn string for database connection
func getDSN(cfg *configs.DBConfig) (string, error) {
	var dsn string
	switch cfg.Driver {
	case "postgres":
		dsn = fmt.Sprintf("host=%s port=%s dbname=%s user=%s password=%s sslmode=%s",
			cfg.Host,
			cfg.Port,
			cfg.Name,
			cfg.User,
			cfg.Password,
			cfg.SSLMode,
		)
	case "mysql":
		dsn = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s",
			cfg.User,
			cfg.Password,
			cfg.Host,
			cfg.Port,
			cfg.Name,
		)
	case "sqlite":
		dsn = fmt.Sprintf("file:%s?cache=%s&mode%s")
	default:
		dsn := os.Getenv("DB_DSN")
		if dsn == "" {
			return dsn, errors.New("supports only postgres, mysql and sqlite")
		}
		return dsn, nil
	}

	return dsn, nil
}
