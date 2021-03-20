package cli_commands

const (
	exitStatusSuccess = 0
	exitStatusError   = 1
)

type serviceContract interface {
	Prepare() error
	CreateMigrationFile(migrationName string) ([]string, error)
	ApplyMigrationsUp() ([]string, error)
	ApplyMigrationsDown() ([]string, error)
	ApplyAllMigrationsDown() ([]string, error)
	RefreshMigrations() ([]string, error)
}

// AbstractCommand
type AbstractCommand struct {
	service serviceContract
}
