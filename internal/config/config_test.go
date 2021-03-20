package config_test

import (
	"github.com/joho/godotenv"
	"github.com/malyg1n/sql-migrator/internal/config"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewConfig(t *testing.T) {
	err := godotenv.Load("../../.env.testing")
	if err != nil {
		t.Fatal(err)
	}
	cfg := config.NewConfig()
	assert.Equal(t, "fake-db", cfg.DbDriver)
	assert.Equal(t, "fake-connection-string", cfg.DbConnectionsString)
	assert.Equal(t, "fake-migrations-folder", cfg.MigrationsPath)
}
