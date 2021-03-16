package tests

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRepository_CreateMigrationsTable(t *testing.T) {
	err := repo.CreateMigrationsTable(GetCreateMigrationsTableSql())
	assert.Nil(t, err)
}

func TestRepository_WorkWithMigrations(t *testing.T) {
	ms, err := repo.GetMigrations()
	assert.Nil(t, err)
	assert.Len(t, ms, 0)

	v, err := repo.GetLatestVersionNumber()
	assert.Nil(t, err)
	assert.Equal(t, 1, v)

	err = repo.ApplyMigrationsUp("create-users", GetCreateUsersTableSql(), v)
	assert.Nil(t, err)

	ms, err = repo.GetMigrations()
	assert.Nil(t, err)
	assert.Len(t, ms, 1)

	v, err = repo.GetLatestVersionNumber()
	assert.Nil(t, err)
	assert.Equal(t, 2, v)

	err = repo.ApplyMigrationsUp("create-users", GetCreateListsTableSql(), v)
	assert.Nil(t, err)

	ms, err = repo.GetMigrations()
	assert.Nil(t, err)
	assert.Len(t, ms, 2)
}
