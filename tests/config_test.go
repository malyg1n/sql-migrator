package tests

import (
	"github.com/malyg1n/sql-migrator/pkg/configs"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewConfig(t *testing.T) {
	cfg := configs.NewConfig()
	assert.Equal(t, "sqlite3", cfg.DB.Driver)
	assert.Equal(t, "test.db", cfg.DB.File)
	assert.Equal(t, "test_migrations", cfg.Main.MigrationsPath)
}
