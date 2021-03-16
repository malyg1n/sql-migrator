package tests

import (
	"database/sql"
	"github.com/joho/godotenv"
	"github.com/malyg1n/sql-migrator/pkg/configs"
	"github.com/malyg1n/sql-migrator/pkg/repositories"
	"github.com/malyg1n/sql-migrator/pkg/services"
	"os"
)

var (
	dbCfg                      *configs.DBConfig
	cfg                        *configs.MainConfig
	repository                 repositories.RepositoryContract
	service                    services.ServiceContract
	db                         *sql.DB
	sqlScriptForMigrationTable = "../sql_scripts/create-migrations-table.sql"
)

// Set up before test
func setUp() {
	godotenv.Load("../.env.testing")
	os.Mkdir("../test_migrations", 0777)
	cfg := configs.NewMainConfig()
	dbCfg := configs.NewDBConfig()
	cfg.MigrationsPath = "../test_migrations"
	dbCfg.File = "../test.db"
	db, _ = sql.Open(dbCfg.Driver, "file:"+dbCfg.File)
	repository = repositories.NewRepository(db)
	service = services.NewService(repository, cfg)
}

// Tear down after tests
func tearDown() {
	os.Remove("../test_migrations")
	os.Remove("../test.db")
}

func GetCreateMigrationsTableSql() string {
	return `

`
}
