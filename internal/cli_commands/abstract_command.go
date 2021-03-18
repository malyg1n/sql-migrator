package cli_commands

import "github.com/malyg1n/sql-migrator/internal/entity"

const (
	exitStatusSuccess = 0
	exitStatusError   = 1
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

// AbstractCommand
type AbstractCommand struct {
	service serviceContract
}
