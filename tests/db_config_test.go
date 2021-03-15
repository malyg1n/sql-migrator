package tests

import (
	"github.com/malyg1n/sql-migrator/pkg/configs"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewDBConfig(t *testing.T) {
	cfg := configs.NewDBConfig()
	assert.Equal(t, "sqlite3", cfg.Driver)
	assert.Equal(t, "test.db", cfg.File)
}
