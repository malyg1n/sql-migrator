package service_test

import (
	"github.com/malyg1n/sql-migrator/internal/entity"
	"github.com/malyg1n/sql-migrator/internal/service"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

type serviceContract interface {
	Prepare() error
	CreateFolder() error
	CreateMigrationFile(migrationName string) ([]string, error)
	ApplyMigrationsUp() ([]string, error)
	ApplyMigrationsDown() ([]string, error)
	ApplyAllMigrationsDown() ([]string, error)
	RefreshMigrations() ([]string, error)
	GetMigrationUpFiles(folder string) ([]string, error)
	FilterMigrations(dbMigrations []*entity.MigrationEntity, files []string) []string
}

type migrationStoreStub struct {
	tableName string
	dbDriver  string
}

const (
	migrationFolder    = "test_migration_folder"
	prepareScriptsPath = "../../prepare"
	timeFormat         = "20060102150405"
)

var (
	srv        *service.Service
	migrations []*entity.MigrationEntity
)

func (store *migrationStoreStub) GetDbDriver() string {
	return store.dbDriver
}

func (store *migrationStoreStub) CreateMigrationsTable(query string) error {
	return nil
}

func (store *migrationStoreStub) GetMigrations() ([]*entity.MigrationEntity, error) {
	return migrations, nil
}

func (store *migrationStoreStub) GetMigrationsByVersion(version uint) ([]*entity.MigrationEntity, error) {
	var migrationsByVersion []*entity.MigrationEntity
	for _, m := range migrations {
		if m.Version == version {
			migrationsByVersion = append(migrationsByVersion, m)
		}
	}

	return migrationsByVersion, nil
}

func (store *migrationStoreStub) GetLatestVersionNumber() (uint, error) {
	var version uint
	for _, m := range migrations {
		if m.Version > version {
			version = m.Version
		}
	}

	return version, nil
}

func (store *migrationStoreStub) ApplyMigrationsUp(migrations []*entity.MigrationEntity) error {
	var newMigrations []*entity.MigrationEntity
	for _, m := range migrations {
		newMigrations = append(newMigrations, m)
	}

	return nil
}

func (store *migrationStoreStub) ApplyMigrationsDown(migrations []*entity.MigrationEntity) error {
	for _, m := range migrations {

	}

	return nil
}

func TestMain(m *testing.M) {
	setUp()
	defer tearDown()
	m.Run()
}

func TestService_CreateFolder(t *testing.T) {
	err := srv.CreateFolder()
	assert.Nil(t, err)
	assert.DirExists(t, migrationFolder)

	info, err := os.Stat(migrationFolder)
	assert.Equal(t, "drwxr--r--", info.Mode().String())
}

func TestService_Prepare(t *testing.T) {
	err := srv.Prepare()
	assert.Nil(t, err)
}

func TestService_ApplyMigrationsUp(t *testing.T) {

}

func setUp() {
	repo := &migrationStoreStub{
		tableName: "test_schema_migrations_service",
		dbDriver:  "sqlite3",
	}
	srv = service.NewService(repo, migrationFolder)
}

func tearDown() {
	os.RemoveAll(migrationFolder)
}
