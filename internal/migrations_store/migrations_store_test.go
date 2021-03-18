package migrations_store_test

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/malyg1n/sql-migrator/internal/entities"
	"github.com/malyg1n/sql-migrator/internal/migrations_store"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
	"testing"
)

type store interface {
	CreateMigrationsTable(query string) error
	GetMigrations() ([]*entities.MigrationEntity, error)
	GetMigrationsByVersion(version uint) ([]*entities.MigrationEntity, error)
	GetLatestVersionNumber() (uint, error)
	ApplyMigrationsUp(migrations []*entities.MigrationEntity) error
	ApplyMigrationsDown(migrations []*entities.MigrationEntity) error
}

const (
	dbDriver         = "postgres"
	dbFilename       = "test_migrations_store.db"
	connectionString = "host=localhost port=6432 dbname=app user=forge password=secrea sslmode=disable"
	tableName        = "test_schema_migrations"
)

var (
	st store
	db *sql.DB
)

func TestMain(m *testing.M) {
	setUp()
	//defer tearDown()
	m.Run()
}

//func TestMigrationsStore_CreateMigrationsTable(t *testing.T) {
//	err := st.CreateMigrationsTable(getCreateMigrationsTableSql(tableName))
//	assert.NoError(t, err)
//
//	err = st.CreateMigrationsTable("fake query")
//	assert.Errorf(t, err, `near "fake": syntax error`)
//}

func TestMigrationsStore_ApplyMigrationsUp(t *testing.T) {
	st.CreateMigrationsTable(getCreateMigrationsTableSql(tableName))
	ms, err := st.GetMigrations()
	assert.NoError(t, err)
	assert.Len(t, ms, 0)

	testCases := []struct {
		name            string
		migration       *entities.MigrationEntity
		countMigrations int
		error           string
	}{
		{
			name:            "valid create users",
			migration:       entities.NewMigrationEntity("create-users", getCreateUsersTableSql(), 1),
			countMigrations: 1,
			error:           "",
		},
		{
			name:            "second create users (not uniq name of migration)",
			migration:       entities.NewMigrationEntity("create-users", "CREATE TABLE new_users( id bigserial not null primary key)", 2),
			countMigrations: 1,
			error:           "UNIQUE constraint failed: test_schema_migrations.migration",
		},
		{
			name:            "third create users (table already exists)",
			migration:       entities.NewMigrationEntity("create-users-new", getCreateUsersTableSql(), 2),
			countMigrations: 1,
			error:           "UNIQUE constraint failed: test_schema_migrations.migration",
		},
		{
			name:            "fake query",
			migration:       entities.NewMigrationEntity("create-users-new", "fake-query", 2),
			countMigrations: 1,
			error:           `near "fake": syntax error`,
		},
		{
			name:            "valid create lists",
			migration:       entities.NewMigrationEntity("create-lists", getCreateListsTableSql(), 2),
			countMigrations: 2,
			error:           "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var ms []*entities.MigrationEntity
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

//func TestMigrationsStore_GetMigrations(t *testing.T) {
//	ms, err := st.GetMigrations()
//
//	assert.NoError(t, err)
//	assert.Len(t, ms, 2)
//}
//
//func TestMigrationsStore_GetMigrationsByVersion(t *testing.T) {
//	ms, err := st.GetMigrationsByVersion(2)
//	assert.NoError(t, err)
//	assert.Len(t, ms, 1)
//
//	ms, err = st.GetMigrationsByVersion(1)
//	assert.NoError(t, err)
//	assert.Len(t, ms, 1)
//
//	ms, err = st.GetMigrationsByVersion(100000)
//	assert.NoError(t, err)
//	assert.Len(t, ms, 0)
//}
//
//func TestMigrationsStore_ApplyMigrationsDown(t *testing.T) {
//	ms, err := st.GetMigrations()
//
//	assert.NoError(t, err)
//	assert.Len(t, ms, 2)
//
//	testCases := []struct {
//		name            string
//		migration      *entities.MigrationEntity
//		countMigrations int
//		error           string
//	}{
//		{
//			name:            "valid drop users",
//			migration:      entities.NewMigrationEntity("create-users", getDropListsTableSql(), 1),
//			countMigrations: 1,
//			error:           "",
//		},
//		{
//			name:            "fake query",
//			migration:      entities.NewMigrationEntity("fake-query", "fake-query", 1),
//			countMigrations: 1,
//			error:          `near "fake": syntax error`,
//		},
//		{
//			name:            "valid drop lists",
//			migration:      entities.NewMigrationEntity("create-lists", getDropListsTableSql(), 1),
//			countMigrations: 0,
//			error:           "",
//		},
//		{
//			name:            "invalid table name",
//			migration:      entities.NewMigrationEntity("invalid table", "DROP TABLE invalid", 1),
//			countMigrations: 0,
//			error:           "no such table: invalid",
//		},
//	}
//
//	for _, tc := range testCases {
//		t.Run(tc.name, func(t *testing.T) {
//			var ms []*entities.MigrationEntity
//			ms = append(ms, tc.migration)
//			err := st.ApplyMigrationsDown(ms)
//			if tc.error == "" {
//				assert.NoError(t, err)
//				ms, _ := st.GetMigrations()
//				assert.Len(t, ms, tc.countMigrations)
//			} else {
//				assert.Errorf(t, err, tc.error)
//			}
//		})
//	}
//
//}

func setUp() {
	db, err := sql.Open(dbDriver, connectionString)
	if err != nil {
		panic(err)
	}
	st = migrations_store.NewStore(db, tableName)
}

func tearDown() {
	db, _ := sql.Open(dbDriver, connectionString)
	db.Exec(`drop table test_schema_migrations;

drop table new_users;

drop table users;
`)
}

// Test queries
func getCreateMigrationsTableSql(tableName string) string {
	return fmt.Sprintf(`
CREATE TABLE IF NOT EXISTS %s
(
     id bigserial not null primary key,
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
     id bigserial not null primary key,
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
     id bigserial not null primary key,
    label varchar (255) not null,
    description  varchar (255) not null unique,
    created_at timestamp DEFAULT CURRENT_TIMESTAMP
)
`
}

func getDropListsTableSql() string {
	return `DROP TABLE lists;`
}
