package commands

import (
	"database/sql"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/malyg1n/sql-migrator/pkg/configs"
)

func setUp() {
	godotenv.Load(".env.testing")
	dbCfg := configs.NewDBConfig()
	dsn := fmt.Sprintf("file:%s?cache=%s&mode=%s", cfg.File, cfg.Cache, cfg.Mode)
	db, err := sql.Open(dbCfg.Driver, fmt.Sprintf(""))
}

func tearDown() {

}

