package sql_migrator

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewConfig(t *testing.T) {
	cfg := NewConfig()
	assert.Equal(t, "sqlite3", cfg.DB.Driver)
	assert.Equal(t, "test.db", cfg.DB.File)
	assert.Equal(t, "test_migrations", cfg.Main.MigrationsPath)
}
