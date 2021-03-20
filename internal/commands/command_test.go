package commands_test

import (
	"github.com/malyg1n/sql-migrator/internal/commands"
	"github.com/stretchr/testify/assert"
	"testing"
)

type serviceStub struct {
}

func (s *serviceStub) Prepare() error {
	return nil
}

func (s *serviceStub) CreateMigrationFiles(migrationName string) ([]string, error) {
	return nil, nil
}

func (s *serviceStub) ApplyMigrationsUp() ([]string, error) {
	return nil, nil
}

func (s *serviceStub) ApplyMigrationsDown() ([]string, error) {
	return nil, nil
}

func (s *serviceStub) ApplyAllMigrationsDown() ([]string, error) {
	return nil, nil
}

func (s *serviceStub) RefreshMigrations() ([]string, error) {
	return nil, nil
}

func TestInitCommand_Run(t *testing.T) {
	cmd := commands.NewInitCommand(&serviceStub{})
	assert.Equal(t, 0, cmd.Run([]string{}))
}

func TestCreateCommand_Run(t *testing.T) {
	cmd := commands.NewCreateCommand(&serviceStub{})
	assert.Equal(t, 1, cmd.Run([]string{}))

	cmd = commands.NewCreateCommand(&serviceStub{})
	assert.Equal(t, 0, cmd.Run([]string{"test"}))
}

func TestDownCommand_Run(t *testing.T) {
	cmd := commands.NewDownCommand(&serviceStub{})
	assert.Equal(t, 0, cmd.Run([]string{}))
}

func TestNewRefreshCommand(t *testing.T) {
	cmd := commands.NewRefreshCommand(&serviceStub{})
	assert.Equal(t, 0, cmd.Run([]string{}))
}

func TestNewCleanCommand(t *testing.T) {
	cmd := commands.NewCleanCommand(&serviceStub{})
	assert.Equal(t, 0, cmd.Run([]string{}))
}
