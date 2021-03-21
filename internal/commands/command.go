package commands

const (
	exitStatusSuccess = iota
	exitStatusError
)

type serviceContract interface {
	Prepare() error
	CreateMigrationFiles(migrationName string) ([]string, error)
	ApplyMigrationsUp() ([]string, error)
	ApplyMigrationsDown() ([]string, error)
	ApplyAllMigrationsDown() ([]string, error)
	RefreshMigrations() ([]string, error)
}

// AbstractCommand is parent command for other commands
type AbstractCommand struct {
	service serviceContract
}
