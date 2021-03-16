package sql_migrator

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
	assert.Equal(t, uint(0), v)

	v++
	err = repo.ApplyMigrationsUp("create-users", GetCreateUsersTableSql(), v)
	assert.Nil(t, err)

	ms, err = repo.GetMigrations()
	assert.Nil(t, err)
	assert.Len(t, ms, 1)

	v, err = repo.GetLatestVersionNumber()
	assert.Nil(t, err)
	assert.Equal(t, uint(1), v)

	v++
	err = repo.ApplyMigrationsUp("create-lists", GetCreateListsTableSql(), v)
	assert.Nil(t, err)

	v, err = repo.GetLatestVersionNumber()
	assert.Nil(t, err)
	assert.Equal(t, uint(2), v)

	ms, err = repo.GetMigrationsByVersion(2)
	assert.Nil(t, err)
	assert.Len(t, ms, 1)

	ms, err = repo.GetMigrationsByVersion(1)
	assert.Nil(t, err)
	assert.Len(t, ms, 1)

	err = repo.ApplyMigrationsDown(2, GetDropUsersTableSql())

	ms, err = repo.GetMigrations()
	assert.Nil(t, err)
	assert.Len(t, ms, 1)

	err = repo.ApplyMigrationsDown(1, GetDropListsTableSql())

	ms, err = repo.GetMigrations()
	assert.Nil(t, err)
	assert.Len(t, ms, 0)
}
