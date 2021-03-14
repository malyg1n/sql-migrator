package commands

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/malyg1n/sqlx-migrator/output"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const (
	timeFormat         = "20060102150405"
	migrationTableName = "schema_migrations"
	exitStatusSuccess  = 0
	exitStatusError    = 1
)

var (
	migrationDir = "migrations"
)

// Migration
type MigrationEntity struct {
	Id        uint      `db:"id"`
	Migration string    `db:"migration"`
	Version   uint      `db:"version"`
	CreatedAt time.Time `db:"created_at"`
}

// AbstractCommand
type AbstractCommand struct {
	db *sql.DB
}

func (c *AbstractCommand) checkMigrationsTable() (bool, error) {
	row, err := c.db.Query(fmt.Sprintf("SELECT COUNT(id) from %s", migrationTableName))
	if err != nil {
		return false, err
	}
	return row != nil, nil
}

func (c *AbstractCommand) createMigrationsTable() error {
	sqlQuery := fmt.Sprintf(`CREATE TABLE %s (
    id bigserial not null primary key,
    migration varchar not null unique,
    version int not null,
    created_at timestamp DEFAULT CURRENT_TIMESTAMP
)`, migrationTableName)
	_, err := c.db.Query(sqlQuery)
	if err != nil {
		return err
	}
	return nil
}

func (c *AbstractCommand) getMigrationsFromBD() ([]*MigrationEntity, error) {
	migrations := make([]*MigrationEntity, 0)
	query := fmt.Sprintf("SELECT * from %s", migrationTableName)
	rows, err := c.db.Query(query)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		me := &MigrationEntity{}
		if err := rows.Scan(&me.Id, &me.Migration, &me.Version, &me.CreatedAt); err != nil {
			return nil, err
		}
		migrations = append(migrations, me)
	}

	return migrations, nil
}

func (c *AbstractCommand) getMigrationUpFiles(folder string) ([]string, error) {
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

func (c *AbstractCommand) filterMigrations(dbMigrations []*MigrationEntity, files []string) []string {
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

func (c *AbstractCommand) applyMigrationsUp(migrationsFiles []string, version uint) error {
	for _, file := range migrationsFiles {
		data, err := ioutil.ReadFile(file)
		if err != nil {
			return err
		}
		migrationQuery := fmt.Sprintf("INSERT INTO %s (migration, version) VALUES ($1, $2)", migrationTableName)
		_, err = c.db.Query(migrationQuery, file, version)
		if err != nil {
			return err
		}

		query := string(data)
		_, err = c.db.Query(query)
		if err != nil {
			return err
		}
		output.ShowMessage(fmt.Sprintf("migrated: %s", file))
	}
	return nil
}

func (c *AbstractCommand) getLatestVersionNumber() (uint, error) {
	var versionNumber uint
	me := &MigrationEntity{}
	query := fmt.Sprintf("SELECT * FROM %s ORDER BY version DESC limit 1", migrationTableName)
	row := c.db.QueryRow(query)
	if row == nil {
		versionNumber = 1
	} else {
		row.Scan(&me.Id, &me.Migration, &me.Version, &me.CreatedAt)
		if me.Version == 0 {
			return 0, errors.New("can not find version number")
		}
		versionNumber = me.Version
	}
	return versionNumber, nil
}

//
//func (c *AbstractCommand) getMigrationsByVersion() []*MigrationEntity {
//
//}
//
//func (c *AbstractCommand) deleteMigrationsByVersion() error {
//
//}
//
//func (c *AbstractCommand) deleteAllMigrations() error {
//
//}
