package services

import "github.com/malyg1n/sql-migrator/entities"

type ServiceContract interface {
	CheckMigrationsTable() error
	CreateMigrationsTable() error
	ApplyMigrationsUp(migrationsFiles []string, version uint) error
	GetMigrations() ([]*entities.MigrationEntity, error)
	GetMigrationsByVersion(version uint) ([]*entities.MigrationEntity, error)
	GetLatestVersionNumber() (uint, error)
	RollbackMigration(entity *entities.MigrationEntity) error
	GetMigrationUpFiles(folder string) ([]string, error)
	FilterMigrations(dbMigrations []*entities.MigrationEntity, files []string) []string
	CheckFolder(dir string) error
}

const (
	timeFormat         = "20060102150405"
	migrationTableName = "schema_migrations"
)
