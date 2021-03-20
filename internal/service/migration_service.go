package service

import (
	"fmt"
	"github.com/malyg1n/sql-migrator/internal/config"
	"github.com/malyg1n/sql-migrator/internal/entity"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
)

type storeContract interface {
	CreateMigrationsTable(query string) error
	GetMigrations() ([]*entity.MigrationEntity, error)
	GetMigrationsByVersion(version uint) ([]*entity.MigrationEntity, error)
	GetLatestVersionNumber() (uint, error)
	ApplyMigrationsUp(migrations []*entity.MigrationEntity) error
	ApplyMigrationsDown(migrations []*entity.MigrationEntity) error
}

type Service struct {
	repo storeContract
	cfg  *config.Config
}

func NewService(repo storeContract, cfg *config.Config) *Service {
	return &Service{
		repo: repo,
		cfg:  cfg,
	}
}

func (s *Service) Prepare() error {
	err := s.createFolder()
	if err != nil {
		return err
	}

	data, err := ioutil.ReadFile(path.Join(s.cfg.PrepareScriptsPath, s.cfg.DbDriver+".sql"))
	if err != nil {
		return err
	}

	return s.repo.CreateMigrationsTable(string(data))
}

func (s *Service) CreateMigrationFile(migrationName string) ([]string, error) {
	var messages []string
	files, err := os.ReadDir(s.cfg.MigrationsPath)
	if err != nil {
		return nil, err
	}

	upFileName := fmt.Sprintf("%s-%s-up.sql", fmt.Sprintf("%05d", (len(files)/2)+1), strings.TrimSpace(migrationName))
	pathName := path.Join(s.cfg.MigrationsPath, upFileName)
	fUp, err := os.Create(pathName)

	if err != nil {
		return nil, err
	}

	messages = append(messages, fmt.Sprintf("created migration %s", pathName))

	downFileName := fmt.Sprintf("%s-%s-down.sql", fmt.Sprintf("%05d", (len(files)/2)+1), strings.TrimSpace(migrationName))
	pathName = path.Join(s.cfg.MigrationsPath, downFileName)
	fDown, err := os.Create(pathName)

	if err != nil {
		return nil, err
	}

	messages = append(messages, fmt.Sprintf("created migration %s", pathName))

	defer func() {
		_ = fUp.Close()
		_ = fDown.Close()
	}()

	return messages, nil
}

func (s *Service) ApplyMigrationsUp() ([]string, error) {
	migrations, err := s.repo.GetMigrations()
	if err != nil {
		return nil, err
	}

	files, err := s.getMigrationUpFiles(s.cfg.MigrationsPath)
	if err != nil {
		return nil, err
	}

	newMigrationsFiles := s.filterMigrations(migrations, files)
	if len(newMigrationsFiles) == 0 {
		return nil, nil
	}

	version, err := s.repo.GetLatestVersionNumber()
	if err != nil {
		return nil, err
	}

	// increase version number
	version++

	var migrated []string
	var newMigrations []*entity.MigrationEntity

	for _, file := range newMigrationsFiles {
		data, err := ioutil.ReadFile(file)
		if err != nil {
			return nil, err
		}
		migrated = append(migrated, fmt.Sprintf("migrated: %s", file))
		newMigrations = append(newMigrations, entity.NewMigrationEntity(file, string(data), version))
	}
	err = s.repo.ApplyMigrationsUp(newMigrations)
	if err != nil {
		return nil, err
	}

	return migrated, nil
}

func (s *Service) ApplyMigrationsDown() ([]string, error) {
	version, err := s.repo.GetLatestVersionNumber()
	if err != nil {
		return nil, err
	}

	migrations, err := s.repo.GetMigrationsByVersion(version)
	if err != nil {
		return nil, err
	}

	var rollback []string
	var backMigrations []*entity.MigrationEntity

	for _, m := range migrations {
		filePath := strings.Replace(m.Migration, "-up.sql", "-down.sql", 1)
		data, err := ioutil.ReadFile(filePath)
		if err != nil {
			return nil, err
		}

		rollback = append(rollback, fmt.Sprintf("rolled back: %s", filePath))
		backMigrations = append(backMigrations, entity.NewMigrationEntity(m.Migration, string(data), m.Version))
	}

	err = s.repo.ApplyMigrationsDown(backMigrations)
	if err != nil {
		return nil, err
	}

	return rollback, err
}

func (s *Service) ApplyAllMigrationsDown() ([]string, error) {
	migrations, err := s.repo.GetMigrations()
	if err != nil {
		return nil, err
	}

	var rollback []string

	var backMigrations []*entity.MigrationEntity

	for _, m := range migrations {
		filePath := strings.Replace(m.Migration, "-up.sql", "-down.sql", 1)
		data, err := ioutil.ReadFile(filePath)
		if err != nil {
			return nil, err
		}

		rollback = append(rollback, fmt.Sprintf("rolled back: %s", filePath))
		backMigrations = append(backMigrations, entity.NewMigrationEntity(m.Migration, string(data), m.Version))
	}

	if err := s.repo.ApplyMigrationsDown(backMigrations); err != nil {
		return nil, err
	}

	return rollback, err
}

func (s *Service) RefreshMigrations() ([]string, error) {
	var messages []string
	rolledBack, err := s.ApplyAllMigrationsDown()
	if err != nil {
		return nil, err
	}

	migrated, err := s.ApplyMigrationsUp()
	if err != nil {
		return nil, err
	}

	messages = append(messages, rolledBack...)
	messages = append(messages, migrated...)

	return messages, err
}

func (s *Service) getMigrationUpFiles(folder string) ([]string, error) {
	files := make([]string, 0)
	err := filepath.Walk(folder, func(path string, info os.FileInfo, err error) error {
		if strings.Contains(path, "-up.sql") {
			files = append(files, path)
		}
		return nil
	})

	if err != nil {
		return nil, err
	}
	return files, nil
}

func (s *Service) filterMigrations(dbMigrations []*entity.MigrationEntity, files []string) []string {
	newFiles := make([]string, 0)
	for _, file := range files {
		found := false
		for _, m := range dbMigrations {
			if m.Migration == file {
				found = true
				break
			}
		}
		if !found {
			newFiles = append(newFiles, file)
		}
	}
	return newFiles
}

func (s *Service) createFolder() error {
	if _, err := os.Stat(s.cfg.MigrationsPath); os.IsExist(err) {
		return os.Mkdir(s.cfg.MigrationsPath, 0764)
	}

	return nil
}
