package migrations_store_test

import (
	"database/sql"
	"github.com/malyg1n/sql-migrator/internal/entities"
	"github.com/malyg1n/sql-migrator/internal/helpers"
	"github.com/malyg1n/sql-migrator/internal/migrations_store"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

type store interface {
	CreateMigrationsTable(query string) error
	GetMigrations() ([]*entities.MigrationEntity, error)
	GetMigrationsByVersion(version uint) ([]*entities.MigrationEntity, error)
	GetLatestVersionNumber() (uint, error)
	ApplyMigrationsUp(migrationName string, dbQuery string, version uint) error
	ApplyMigrationsDown(migrationId uint, dbQuery string) error
}

const (
	dbDriver         = "sqlite3"
	dbFilename       = "test_migrations_store.db"
	connectionString = "file:" + dbFilename
	tableName        = "test_schema_migrations"
)

var (
	st store
	db *sql.DB
)

func TestMain(m *testing.M) {
	setUp()
	m.Run()
	tearDown()
}

func TestMigrationsStore_CreateMigrationsTable(t *testing.T) {
	err := st.CreateMigrationsTable(helpers.GetCreateMigrationsTableSql(tableName))
	assert.NoError(t, err)

	err = st.CreateMigrationsTable("fake query")
	assert.Errorf(t, err, `near "fake": syntax error`)
}

func TestMigrationsStore_ApplyMigrationsUp(t *testing.T) {
	ms, err := st.GetMigrations()
	assert.NoError(t, err)
	assert.Len(t, ms, 0)

	testCases := []struct {
		name            string
		migrationName   string
		query           string
		version         uint
		countMigrations int
		error           string
	}{
		{
			name:            "valid create users",
			migrationName:   "create-users",
			query:           helpers.GetCreateUsersTableSql(),
			version:         1,
			countMigrations: 1,
			error:           "",
		},
		{
			name:            "second create users (not uniq name of migration)",
			migrationName:   "create-users",
			query:           helpers.GetCreateUsersTableSql(),
			version:         1,
			countMigrations: 1,
			error:           "UNIQUE constraint failed: test_schema_migrations.migration",
		},
		{
			name:            "fake query",
			migrationName:   "fake-query",
			query:           helpers.GetCreateUsersTableSql(),
			version:         1,
			countMigrations: 1,
			error:           "UNIQUE constraint failed: test_schema_migrations.migration",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := st.ApplyMigrationsUp(tc.migrationName, tc.query, tc.version)
			if tc.error == "" {
				assert.NoError(t, err)
				ms, _ := st.GetMigrations()
				assert.Len(t, ms, 1)
			} else {
				assert.Errorf(t, err, tc.error)
			}
		})
	}
}

func setUp() {
	db, err := sql.Open(dbDriver, connectionString)
	if err != nil {
		panic(err)
	}
	st = migrations_store.NewStore(db, tableName)
}

func tearDown() {
	os.Remove(dbFilename)
}
