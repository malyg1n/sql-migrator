package repositories

import (
	"github.com/malyg1n/sql-migrator/pkg/entities"
)

const (
	migrationTableName = "schema_migrations"
)

type RepositoryContract interface {
	CheckOrCreateMigrationsTable() error
	ApplyMigrationsUp(migrationName string, content string, version uint) error
	GetMigrations() ([]*entities.MigrationEntity, error)
	GetMigrationsByVersion(version uint) ([]*entities.MigrationEntity, error)
	GetLatestVersionNumber() (uint, error)
	ApplyMigrationsDown(migrationId uint, content string) error
}
