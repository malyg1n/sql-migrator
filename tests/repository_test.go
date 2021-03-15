package tests

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCheckOrCreateMigrationsTable(t *testing.T) {
	err := repository.CheckOrCreateMigrationsTable(sqlScriptForMigrationTable)
	assert.Nil(t, err)
}
