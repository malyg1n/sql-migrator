package services

import "github.com/malyg1n/sql-migrator/pkg/entities"

type ServiceContract interface {
	Prepare() error
	CreateMigrationFile(migrationName string) ([]string, error)
	ApplyMigrationsUp() ([]string, error)
	ApplyMigrationsDown() ([]string, error)
	ApplyAllMigrationsDown() ([]string, error)
	RefreshMigrations() ([]string, error)
	GetMigrationUpFiles(folder string) ([]string, error)
	FilterMigrations(dbMigrations []*entities.MigrationEntity, files []string) []string
}

const (
	timeFormat = "20060102150405"
)
