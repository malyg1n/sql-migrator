package services

import (
	"errors"
	"fmt"
	"github.com/malyg1n/sql-migrator/entities"
	"github.com/malyg1n/sql-migrator/output"
	"github.com/malyg1n/sql-migrator/repositories"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

type Service struct {
	repo repositories.RepositoryContract
}

var (
	migrationDir = "migrations"
)

func (s *Service) CheckMigrationsTable() error {
	return s.repo.CheckMigrationsTable()
}

func (s *Service) CreateMigrationsTable() error {
	return s.repo.CreateMigrationsTable()
}

func (s *Service) ApplyMigrationsUp(migrationsFiles []string, version uint) error {
	for _, file := range migrationsFiles {
		data, err := ioutil.ReadFile(file)
		if err != nil {
			return err
		}
		err = s.repo.ApplyMigrationsUp(file, string(data), version)
		if err != nil {
			return err
		}
		output.ShowMessage(fmt.Sprintf("migrated: %s", file))
	}
	return nil
}

func (s *Service) GetMigrations() ([]*entities.MigrationEntity, error) {
	return s.repo.GetMigrations()
}

func (s *Service) GetMigrationsByVersion(version uint) ([]*entities.MigrationEntity, error) {
	return s.repo.GetMigrationsByVersion(version)
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

func (s *Service) FilterMigrations(dbMigrations []*entities.MigrationEntity, files []string) []string {
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

func (s *Service) GetLatestVersionNumber() (uint, error) {
	var versionNumber uint
	me := &entities.MigrationEntity{}
	query := fmt.Sprintf("SELECT * FROM %s ORDER BY version DESC limit 1", migrationTableName)
	row := c.db.QueryRow(query)
	if row == nil {
		versionNumber = 1
	} else {
		row.Scan(&me.Id, &me.Migration, &me.Version, &me.CreatedAt)
		if me.Version == 0 {
			return 1, nil
		}
		versionNumber = me.Version
	}
	return versionNumber, nil
}

func (s *Service) RollbackMigration(entity *entities.MigrationEntity) error {
	filePath := strings.Replace(entity.Migration, "-up.sql", "-down.sql", 1)
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}

	downQuery := string(data)
	_, err = c.db.Query(downQuery)
	if err != nil {
		return err
	}

	query := fmt.Sprintf("DELETE FROM %s WHERE id=%d", migrationTableName, entity.Id)
	_, err = c.db.Query(query)
	if err != nil {
		return err
	}

	output.ShowWarning(fmt.Sprintf("rolled back: %s", entity.Migration))
	return nil
}

func (s *Service) CheckFolder(dir string) error {
	_, err := os.Stat(dir)
	if !os.IsNotExist(err) {
		return nil
	}
	return errors.New(fmt.Sprintf(" no such file or directory: %s", dir))
}
