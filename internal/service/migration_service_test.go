package service_test

import (
	"fmt"
	"github.com/malyg1n/sql-migrator/internal/config"
	"github.com/malyg1n/sql-migrator/internal/entity"
	"github.com/malyg1n/sql-migrator/internal/service"
	"github.com/stretchr/testify/assert"
	"os"
	"path"
	"sort"
	"sync"
	"testing"
)

type serviceContract interface {
	Prepare() error
	CreateMigrationFile(migrationName string) ([]string, error)
	ApplyMigrationsUp() ([]string, error)
	ApplyMigrationsDown() ([]string, error)
	ApplyAllMigrationsDown() ([]string, error)
	RefreshMigrations() ([]string, error)
}

type migrationStoreStub struct {
	tableName      string
	dbDriver       string
	mx             sync.Mutex
	fakeMigrations map[string]*entity.MigrationEntity
}

const (
	migrationFolder     = "test_migration_folder"
	firstMigrationName  = "000001-first"
	secondMigrationName = "000002-second"
)

var (
	srv serviceContract
)

func (store *migrationStoreStub) GetDbDriver() string {
	return store.dbDriver
}

func (store *migrationStoreStub) CreateMigrationsTable(query string) error {
	return nil
}

func (store *migrationStoreStub) GetMigrations() ([]*entity.MigrationEntity, error) {
	store.mx.Lock()
	defer store.mx.Unlock()
	var migrations []*entity.MigrationEntity
	keys := make([]string, 0)
	for k := range store.fakeMigrations {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for i := len(keys) - 1; i >= 0; i-- {
		migrations = append(migrations, store.fakeMigrations[keys[i]])
	}
	return migrations, nil
}

func (store *migrationStoreStub) GetMigrationsByVersion(version uint) ([]*entity.MigrationEntity, error) {
	store.mx.Lock()
	defer store.mx.Unlock()
	var migrationsByVersion []*entity.MigrationEntity
	keys := make([]string, 0)
	for k, m := range store.fakeMigrations {
		if m.Version == version {
			keys = append(keys, k)
		}

	}
	sort.Strings(keys)
	for i := len(keys) - 1; i >= 0; i-- {
		migrationsByVersion = append(migrationsByVersion, store.fakeMigrations[keys[i]])
	}

	return migrationsByVersion, nil
}

func (store *migrationStoreStub) GetLatestVersionNumber() (uint, error) {
	var version uint
	for _, m := range store.fakeMigrations {
		if m.Version > version {
			version = m.Version
		}
	}

	return version, nil
}

func (store *migrationStoreStub) ApplyMigrationsUp(migrations []*entity.MigrationEntity) error {
	for _, m := range migrations {
		store.mx.Lock()
		store.fakeMigrations[m.Migration] = m
		store.mx.Unlock()
	}

	return nil
}

func (store *migrationStoreStub) ApplyMigrationsDown(migrations []*entity.MigrationEntity) error {
	for _, m := range migrations {
		if _, ok := store.fakeMigrations[m.Migration]; ok {
			store.mx.Lock()
			delete(store.fakeMigrations, m.Migration)
			store.mx.Unlock()
		}
	}
	return nil
}

func TestMain(m *testing.M) {
	setUp()
	code := m.Run()
	tearDown()
	os.Exit(code)
}

func TestService_Prepare(t *testing.T) {
	err := srv.Prepare()
	assert.Nil(t, err)
}

func TestService_CreateMigrationFile(t *testing.T) {
	testMigrationName := "test-migration"
	messages, err := srv.CreateMigrationFile(testMigrationName)
	pathNameUp := path.Join(migrationFolder, fmt.Sprintf("%s-%s-up.sql", "00001", testMigrationName))
	pathNameDown := path.Join(migrationFolder, fmt.Sprintf("%s-%s-down.sql", "00001", testMigrationName))
	assert.Nil(t, err)
	assert.FileExists(t, pathNameUp)
	assert.Len(t, messages, 2)
	assert.Equal(t, fmt.Sprintf("created migration %s", pathNameUp), messages[0])
	assert.Equal(t, fmt.Sprintf("created migration %s", pathNameDown), messages[1])
	os.Remove(pathNameUp)
	os.Remove(pathNameDown)
}

func TestService_ApplyMigrationsUp(t *testing.T) {
	testCases := []struct {
		name          string
		migrationUp   string
		migrationDown string
		error         string
		message       string
	}{
		{
			name:          "empty migration",
			migrationUp:   "",
			migrationDown: "",
			error:         "",
			message:       "",
		},
		{
			name:          "valid migration",
			migrationUp:   path.Join(migrationFolder, firstMigrationName+"-up.sql"),
			migrationDown: path.Join(migrationFolder, firstMigrationName+"-down.sql"),
			error:         "",
			message:       "migrated: " + path.Join(migrationFolder, firstMigrationName+"-up.sql"),
		},
		{
			name:          "valid migration 2",
			migrationUp:   path.Join(migrationFolder, secondMigrationName+"-up.sql"),
			migrationDown: path.Join(migrationFolder, secondMigrationName+"-down.sql"),
			error:         "",
			message:       "migrated: " + path.Join(migrationFolder, secondMigrationName+"-up.sql"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.migrationUp == "" {
				messages, err := srv.ApplyMigrationsUp()
				assert.Nil(t, err)
				assert.Nil(t, messages)
			} else {
				createMigrationFiles(t, tc.migrationUp)
				createMigrationFiles(t, tc.migrationDown)
				messages, err := srv.ApplyMigrationsUp()
				assert.Equal(t, tc.message, messages[0])
				assert.Nil(t, err)
			}
		})
	}
}

func TestService_ApplyMigrationsDown(t *testing.T) {
	messages, err := srv.ApplyMigrationsDown()
	assert.Nil(t, err)
	assert.Equal(t, "rolled back: "+path.Join(migrationFolder, secondMigrationName+"-down.sql"), messages[0])

	messages, err = srv.ApplyMigrationsDown()
	assert.Nil(t, err)
	assert.Equal(t, "rolled back: "+path.Join(migrationFolder, firstMigrationName+"-down.sql"), messages[0])
}

func TestService_RefreshMigrations(t *testing.T) {
	messages, err := srv.RefreshMigrations()
	assert.Nil(t, err)
	assert.Len(t, messages, 2)
	assert.Equal(t, "migrated: "+path.Join(migrationFolder, firstMigrationName+"-up.sql"), messages[0])
	assert.Equal(t, "migrated: "+path.Join(migrationFolder, secondMigrationName+"-up.sql"), messages[1])

	_, err = srv.ApplyMigrationsUp()
	if err != nil {
		t.Fatal(err)
	}

	messages, err = srv.RefreshMigrations()
	assert.Nil(t, err)
	assert.Len(t, messages, 4)
	assert.Equal(t, "rolled back: "+path.Join(migrationFolder, secondMigrationName+"-down.sql"), messages[0])
	assert.Equal(t, "rolled back: "+path.Join(migrationFolder, firstMigrationName+"-down.sql"), messages[1])
	assert.Equal(t, "migrated: "+path.Join(migrationFolder, firstMigrationName+"-up.sql"), messages[2])
	assert.Equal(t, "migrated: "+path.Join(migrationFolder, secondMigrationName+"-up.sql"), messages[3])

}

func TestService_ApplyAllMigrationsDown(t *testing.T) {
	messages, err := srv.ApplyAllMigrationsDown()
	assert.Nil(t, err)
	assert.Len(t, messages, 2)
	assert.Equal(t, "rolled back: "+path.Join(migrationFolder, secondMigrationName+"-down.sql"), messages[0])
	assert.Equal(t, "rolled back: "+path.Join(migrationFolder, firstMigrationName+"-down.sql"), messages[1])
}

func setUp() {
	cfg := config.NewConfig()
	cfg.DbDriver = "sqlite3"
	cfg.MigrationsPath = "test_migration_folder"
	repo := &migrationStoreStub{
		tableName:      "test_schema_migrations_service",
		fakeMigrations: make(map[string]*entity.MigrationEntity),
	}
	srv = service.NewService(repo, cfg)
}

func tearDown() {
	os.RemoveAll(migrationFolder)
}

func createMigrationFiles(t *testing.T, filename string) {
	_, err := os.Create(filename)
	if err != nil {
		t.Fatalf("can't create test_migration_folder with error: %s", err)
	}
}
