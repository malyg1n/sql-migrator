package main

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/malyg1n/sql-migrator/sql-migrator"
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
		sql_migrator.ShowError("error loading environment file")
		os.Exit(1)
	}

	cfg := sql_migrator.NewConfig()
	db, err := InitDB(cfg.DB)
	defer db.Close()

	if err != nil {
		sql_migrator.ShowError(err.Error())
		os.Exit(1)
	}

	repo := sql_migrator.NewRepository(db)
	service := sql_migrator.NewService(repo, cfg)

	status, err := InitCommands(service)
	if err != nil {
		sql_migrator.ShowError(err.Error())
		os.Exit(1)
	}

	os.Exit(status)
}

// Init list of commands
func InitCommands(service sql_migrator.ServiceContract) (int, error) {
	c := cli.NewCLI("migrator", "0.0.5")
	c.Args = os.Args[1:]
	c.Commands = map[string]cli.CommandFactory{
		"init": func() (cli.Command, error) {
			return sql_migrator.NewInitCommand(service), nil
		},
		"create": func() (cli.Command, error) {
			return sql_migrator.NewCreateCommand(service), nil
		},
		"up": func() (cli.Command, error) {
			return sql_migrator.NewUpCommand(service), nil
		},
		"down": func() (cli.Command, error) {
			return sql_migrator.NewDownCommand(service), nil
		},
		"refresh": func() (cli.Command, error) {
			return sql_migrator.NewRefreshCommand(service), nil
		},
		"clean": func() (cli.Command, error) {
			return sql_migrator.NewCleanCommand(service), nil
		},
	}

	return c.Run()
}

// Init connect to database
func InitDB(cfg *sql_migrator.DBConfig) (sql_migrator.DBContract, error) {
	dsn, err := GetDSN(cfg)
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
func GetDSN(cfg *sql_migrator.DBConfig) (string, error) {
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
	case "sqlite3":
		dsn = fmt.Sprintf("file:%s", cfg.File)
	default:
		dsn := os.Getenv("DB_DSN")
		if dsn == "" {
			return dsn, errors.New("you must specify the dsn for the database")
		}
		return dsn, nil
	}

	return dsn, nil
}
