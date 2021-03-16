package main

import (
	"github.com/joho/godotenv"
	"github.com/malyg1n/sql-migrator/pkg/configs"
	"github.com/malyg1n/sql-migrator/pkg/repositories"
	"github.com/malyg1n/sql-migrator/pkg/services"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	godotenv.Load(".env.testing")
	code := m.Run()
	os.Exit(code)
}

func TestGetDSN(t *testing.T) {
	cfg := configs.NewDBConfig()

	testCases := []struct {
		name   string
		driver string
		error  interface{}
		answer string
	}{
		{
			name:   "invalid driver",
			driver: "invalid",
			error:  true,
			answer: "",
		},
		{
			name:   "valid driver",
			driver: "sqlite3",
			error:  false,
			answer: "file:test.db",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			cfg.Driver = tc.driver
			dsn, err := GetDSN(cfg)
			assert.Equal(t, tc.answer, dsn)
			if tc.error == true {
				assert.EqualError(t, err, "you must specify the dsn for the database")
			}
		})
	}
}

func TestInitDB(t *testing.T) {
	cfg := configs.NewConfig()
	db, err := InitDB(cfg.DB)
	defer db.Close()
	assert.Nil(t, err)
}

func TestInitCommands(t *testing.T) {
	dbCfg := configs.NewDBConfig()
	db, err := InitDB(dbCfg)
	defer db.Close()

	assert.Nil(t, err)
	cfg := configs.NewConfig()

	repo := repositories.NewRepository(db)
	service := services.NewService(repo, cfg)

	_, err = InitCommands(service)
	assert.Nil(t, err)
}
