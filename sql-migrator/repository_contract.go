package sql_migrator

const (
	migrationTableName = "schema_migrations"
)

type RepositoryContract interface {
	CreateMigrationsTable(query string) error
	ApplyMigrationsUp(migrationName string, content string, version uint) error
	GetMigrations() ([]*MigrationEntity, error)
	GetMigrationsByVersion(version uint) ([]*MigrationEntity, error)
	GetLatestVersionNumber() (uint, error)
	ApplyMigrationsDown(migrationId uint, content string) error
}
