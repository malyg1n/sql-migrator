package store_test

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/malyg1n/sql-migrator/internal/entity"
	"github.com/malyg1n/sql-migrator/internal/store"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

type storeContract interface {
	CreateMigrationsTable(query string) error
	GetMigrations() ([]*entity.MigrationEntity, error)
	GetMigrationsByVersion(version uint) ([]*entity.MigrationEntity, error)
	GetLatestVersionNumber() (uint, error)
	ApplyMigrationsUp(migrations []*entity.MigrationEntity) error
	ApplyMigrationsDown(migrations []*entity.MigrationEntity) error
}

const (
	dbDriver         = "sqlite3"
	dbFilename       = "test_migrations_store.db"
	connectionString = "file:" + dbFilename
	tableName        = "test_schema_migrations"
)

var (
	st storeContract
)

func TestMain(m *testing.M) {
	setUp()
	code := m.Run()
	tearDown()
	os.Exit(code)
}

func TestMigrationsStore_CreateMigrationsTable(t *testing.T) {
	err := st.CreateMigrationsTable(getCreateMigrationsTableSql(tableName))
	assert.NoError(t, err)

	err = st.CreateMigrationsTable("fake query")
	assert.Errorf(t, err, `near "fake": syntax error`)
}

func TestMigrationsStore_ApplyMigrationsUp(t *testing.T) {
	err := st.CreateMigrationsTable(getCreateMigrationsTableSql(tableName))
	if err != nil {
		t.Fatal(err)
	}

	ms, err := st.GetMigrations()
	if err != nil {
		t.Fatal(err)
	}
	assert.Len(t, ms, 0)

	testCases := []struct {
		name            string
		migration       *entity.MigrationEntity
		countMigrations int
		error           string
	}{
		{
			name:            "valid create users",
			migration:       entity.NewMigrationEntity("create-users", getCreateUsersTableSql(), 1),
			countMigrations: 1,
			error:           "",
		},
		{
			name:            "third create users (table already exists)",
			migration:       entity.NewMigrationEntity("create-users-new", getCreateUsersTableSql(), 2),
			countMigrations: 1,
			error:           "UNIQUE constraint failed: test_schema_migrations.migration",
		},
		{
			name:            "fake query",
			migration:       entity.NewMigrationEntity("fake-query", "fake-query", 2),
			countMigrations: 1,
			error:           `near "fake": syntax error`,
		},
		{
			name:            "valid create lists",
			migration:       entity.NewMigrationEntity("create-lists", getCreateListsTableSql(), 2),
			countMigrations: 2,
			error:           "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var ms []*entity.MigrationEntity
			ms = append(ms, tc.migration)
			err := st.ApplyMigrationsUp(ms)
			if tc.error == "" {
				assert.NoError(t, err)
				ms, _ := st.GetMigrations()
				assert.Len(t, ms, tc.countMigrations)
			} else {
				assert.Errorf(t, err, tc.error)
			}
		})
	}
}

func TestMigrationsStore_GetMigrations(t *testing.T) {
	ms, err := st.GetMigrations()

	assert.NoError(t, err)
	assert.Len(t, ms, 2)
}

func TestMigrationsStore_GetMigrationsByVersion(t *testing.T) {
	ms, err := st.GetMigrationsByVersion(2)
	assert.NoError(t, err)
	assert.Len(t, ms, 1)

	ms, err = st.GetMigrationsByVersion(1)
	assert.NoError(t, err)
	assert.Len(t, ms, 1)

	ms, err = st.GetMigrationsByVersion(100000)
	assert.NoError(t, err)
	assert.Len(t, ms, 0)
}

func TestMigrationsStore_ApplyMigrationsDown(t *testing.T) {
	ms, err := st.GetMigrations()

	assert.NoError(t, err)
	assert.Len(t, ms, 2)

	testCases := []struct {
		name            string
		migration       *entity.MigrationEntity
		countMigrations int
		error           string
	}{
		{
			name:            "valid drop users",
			migration:       entity.NewMigrationEntity("create-users", getDropUsersTableSql(), 1),
			countMigrations: 1,
			error:           "",
		},
		{
			name:            "fake query",
			migration:       entity.NewMigrationEntity("fake-query", "fake-query", 1),
			countMigrations: 1,
			error:           `near "fake": syntax error`,
		},
		{
			name:            "valid drop lists",
			migration:       entity.NewMigrationEntity("create-lists", getDropListsTableSql(), 1),
			countMigrations: 0,
			error:           "",
		},
		{
			name:            "invalid table name",
			migration:       entity.NewMigrationEntity("invalid table", "DROP TABLE invalid", 1),
			countMigrations: 0,
			error:           "no such table: invalid",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var ms []*entity.MigrationEntity
			ms = append(ms, tc.migration)
			err := st.ApplyMigrationsDown(ms)
			if tc.error == "" {
				assert.NoError(t, err)
				ms, _ := st.GetMigrations()
				assert.Len(t, ms, tc.countMigrations)
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
	st = store.NewStore(db, tableName)
}

func tearDown() {
	os.Remove(dbFilename)
}

// Test queries
func getCreateMigrationsTableSql(tableName string) string {
	return fmt.Sprintf(`
CREATE TABLE IF NOT EXISTS %s
(
	id integer not null primary key autoincrement,
    migration varchar(255) not null unique,
    version int not null,
    created_at timestamp DEFAULT CURRENT_TIMESTAMP
);
`, tableName)
}

func getCreateUsersTableSql() string {
	return `
CREATE TABLE users
(
    id integer not null primary key autoincrement,
    name varchar (255) not null,
    email  varchar (255) not null unique,
    created_at timestamp DEFAULT CURRENT_TIMESTAMP
)
`
}

func getDropUsersTableSql() string {
	return `DROP TABLE users;`
}

func getCreateListsTableSql() string {
	return `
CREATE TABLE lists
(
    id integer not null primary key autoincrement,
    label varchar (255) not null,
    description  varchar (255) not null unique,
    created_at timestamp DEFAULT CURRENT_TIMESTAMP
)
`
}

func getDropListsTableSql() string {
	return `DROP TABLE lists;`
}
