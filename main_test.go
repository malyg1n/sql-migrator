package main

import (
	"github.com/joho/godotenv"
	"github.com/malyg1n/sql-migrator/configs"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestGetDSN(t *testing.T) {
	godotenv.Load(".env.testing")
	cfg := &configs.DBConfig{
		Driver: os.Getenv("DB_DRIVER"),
		File:   os.Getenv("DB_FILE"),
		Cache:  os.Getenv("DB_CACHE"),
		Mode:   os.Getenv("DB_MODE"),
	}

	assert.Equal(t, "sqlite3", cfg.Driver)
	assert.Equal(t, "test.db", cfg.File)
	assert.Equal(t, "shared", cfg.Cache)
	assert.Equal(t, "memory", cfg.Mode)

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
			answer: "file:test.db?cache=shared&mode=memory",
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
