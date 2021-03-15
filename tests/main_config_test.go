package tests

import (
	"github.com/malyg1n/sql-migrator/pkg/configs"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewMainConfig(t *testing.T) {
	cfg := configs.NewMainConfig()
	assert.Equal(t, "test_migrations", cfg.MigrationsPath)
}
