package migration_service

import (
	"fmt"
	"github.com/malyg1n/sql-migrator/internal"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"
)

type Service struct {
	repo           internal.RepositoryContract
	migrationsPath string
}

func NewService(repo internal.RepositoryContract, migrationsPath string) *Service {
	return &Service{
		repo:           repo,
		migrationsPath: migrationsPath,
	}
}

func (s *Service) Prepare() error {
	err := s.createFolder()
	if err != nil {
		return err
	}

	data, err := ioutil.ReadFile(s.cfg.Main.PrepareScriptsPath + s.cfg.DB.Driver + ".sql")
	if err != nil {
		return err
	}

	return s.repo.CreateMigrationsTable(string(data))
}

func (s *Service) createFolder() error {
	return os.Mkdir(s.cfg.Main.MigrationsPath, 0764)
}

func (s *Service) CreateMigrationFile(migrationName string) ([]string, error) {
	var messages []string
	upFileName := fmt.Sprintf("%s-%s-up.sql", time.Now().Format(internal.timeFormat), strings.TrimSpace(migrationName))
	pathName := path.Join(s.cfg.Main.MigrationsPath, upFileName)
	fUp, err := os.Create(pathName)

	if err != nil {
		return nil, err
	}

	messages = append(messages, fmt.Sprintf("created migration %s", pathName))

	downFileName := fmt.Sprintf("%s-%s-down.sql", time.Now().Format(internal.timeFormat), strings.TrimSpace(migrationName))
	pathName = path.Join(s.cfg.Main.MigrationsPath, downFileName)
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

	files, err := s.GetMigrationUpFiles(s.cfg.Main.MigrationsPath)
	if err != nil {
		return nil, err
	}

	newMigrationsFiles := s.FilterMigrations(migrations, files)
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

	for _, file := range newMigrationsFiles {
		data, err := ioutil.ReadFile(file)
		if err != nil {
			return nil, err
		}
		err = s.repo.ApplyMigrationsUp(file, string(data), version)
		if err != nil {
			return nil, err
		}
		migrated = append(migrated, file)
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

	for _, m := range migrations {
		filePath := strings.Replace(m.Migration, "-up.sql", "-down.sql", 1)
		data, err := ioutil.ReadFile(filePath)
		if err != nil {
			return nil, err
		}

		if err := s.repo.ApplyMigrationsDown(m.Id, string(data)); err != nil {
			return nil, err
		}

		rollback = append(rollback, filePath)
	}

	return rollback, err
}

func (s *Service) ApplyAllMigrationsDown() ([]string, error) {
	migrations, err := s.repo.GetMigrations()
	if err != nil {
		return nil, err
	}

	var rollback []string

	for _, m := range migrations {
		filePath := strings.Replace(m.Migration, "-up.sql", "-down.sql", 1)
		data, err := ioutil.ReadFile(filePath)
		if err != nil {
			return nil, err
		}

		if err := s.repo.ApplyMigrationsDown(m.Id, string(data)); err != nil {
			return nil, err
		}

		rollback = append(rollback, filePath)
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

	for _, rb := range rolledBack {
		messages = append(messages, fmt.Sprintf("rolled back: %s", rb))
	}

	for _, m := range migrated {
		messages = append(messages, fmt.Sprintf("migrated: %s", m))
	}

	return messages, err
}

func (s *Service) GetMigrationUpFiles(folder string) ([]string, error) {
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

func (s *Service) FilterMigrations(dbMigrations []*MigrationEntity, files []string) []string {
	newFiles := make([]string, 0)
	for _, file := range files {
		found := false
		for _, m := range dbMigrations {
			if m.Migration == file {
				found = true
				break
			}
		}
		if found == false {
			newFiles = append(newFiles, file)
		}
	}
	return newFiles
}