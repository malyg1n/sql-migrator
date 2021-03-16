package tests

import (
	"database/sql"
	"github.com/joho/godotenv"
	"github.com/malyg1n/sql-migrator/pkg/configs"
	"github.com/malyg1n/sql-migrator/pkg/database"
	"github.com/malyg1n/sql-migrator/pkg/repositories"
	"github.com/malyg1n/sql-migrator/pkg/services"
	_ "github.com/mattn/go-sqlite3"
	"os"
	"testing"
)

var (
	cfg     *configs.Config
	repo    repositories.RepositoryContract
	service services.ServiceContract
	db      database.DBContract
)

func TestMain(m *testing.M) {
	setUp()
	code := m.Run()
	tearDown()
	os.Exit(code)
}

// Set up before test
func setUp() {
	godotenv.Load("../.env.testing")
	os.Mkdir("../test_migrations", 0777)
	cfg := configs.NewConfig()
	cfg.Main.MigrationsPath = "../test_migrations"
	cfg.DB.File = "../test.db"
	cfg.Main.PrepareScriptsPath = ".." + cfg.Main.PrepareScriptsPath
	db, err := sql.Open(cfg.DB.Driver, "file:"+cfg.DB.File)
	if err != nil {
		panic(err)
	}
	repo = repositories.NewRepository(db)
	service = services.NewService(repo, cfg)
}

// Tear down after tests
func tearDown() {
	os.Remove("../test_migrations")
	os.Remove("../test.db")
}
