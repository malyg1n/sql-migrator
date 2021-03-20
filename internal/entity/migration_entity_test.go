package entity_test

import (
	"github.com/malyg1n/sql-migrator/internal/entity"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewMigrationEntity(t *testing.T) {
	me := entity.NewMigrationEntity("some-migration-up", "some sql query", 1)
	assert.IsType(t, &entity.MigrationEntity{}, me)
}
