package main

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/malyg1n/sql-migrator/pkg/commands"
	"github.com/malyg1n/sql-migrator/pkg/configs"
	"github.com/malyg1n/sql-migrator/pkg/output"
	"github.com/malyg1n/sql-migrator/pkg/repositories"
	"github.com/malyg1n/sql-migrator/pkg/services"
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

	dbCfg := configs.NewDBConfig()
	db, err := InitDB(dbCfg)

	if err != nil {
		output.ShowError(err.Error())
		os.Exit(1)
	}
	cfg := configs.NewMainConfig()

	repo := repositories.NewRepository(db)
	service := services.NewService(repo, cfg)

	err = InitTable(dbCfg, db)
	if err != nil {
		output.ShowError(err.Error())
		os.Exit(1)
	}

	status, err := InitCommands(service)
	if err != nil {
		output.ShowError(err.Error())
		os.Exit(1)
	}

	os.Exit(status)
}

// Init list of commands
func InitCommands(service services.ServiceContract) (int, error) {
	c := cli.NewCLI("migrator", "0.0.5")
	c.Args = os.Args[1:]
	c.Commands = map[string]cli.CommandFactory{
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

	return c.Run()
}

// Init connect to database
func InitDB(cfg *configs.DBConfig) (*sql.DB, error) {
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
func GetDSN(cfg *configs.DBConfig) (string, error) {
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

// Initialize migrations table
func InitTable(cfg *configs.DBConfig, db *sql.DB) error {
	var sql string
	switch cfg.Driver {
	case "postgres":
		{
			sql = `
CREATE TABLE IF NOT EXISTS schema_migrations
(
    id bigserial not null primary key,
    migration varchar(255) not null unique,
    version int not null,
    created_at timestamp DEFAULT CURRENT_TIMESTAMP
);
`
		}
	default:
		sql = `
CREATE TABLE IF NOT EXISTS schema_migrations
(
    id integer not null primary key auto_increment,
    migration varchar(255) not null unique,
    version int not null,
    created_at timestamp DEFAULT CURRENT_TIMESTAMP
);
`
	}

	_, err := db.Exec(sql)

	return err
}
