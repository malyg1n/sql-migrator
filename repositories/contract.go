package repositories

import (
	"github.com/malyg1n/sql-migrator/entities"
)

type RepositoryContract interface {
	CheckMigrationsTable() error
	CreateMigrationsTable() error
	ApplyMigrationsUp(migrationName string, content string, version uint) error
	GetMigrations() ([]*entities.MigrationEntity, error)
	GetMigrationsByVersion(version uint) ([]*entities.MigrationEntity, error)
	GetLatestVersionNumber() (uint, error)
	RollbackMigration(entity *entities.MigrationEntity) error
}
