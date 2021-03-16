package sql_migrator

import (
	"database/sql"
	"github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3"
	"os"
	"testing"
)

var (
	cfg     *Config
	repo    RepositoryContract
	service ServiceContract
	db      DBContract
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
	cfg := NewConfig()
	cfg.Main.MigrationsPath = "../test_migrations"
	cfg.DB.File = "../test.db"
	cfg.Main.PrepareScriptsPath = ".." + cfg.Main.PrepareScriptsPath
	db, err := sql.Open(cfg.DB.Driver, "file:"+cfg.DB.File)
	if err != nil {
		panic(err)
	}
	repo = NewRepository(db)
	service = NewService(repo, cfg)
}

// Tear down after tests
func tearDown() {
	os.Remove("../test_migrations")
	os.Remove("../test.db")
}
