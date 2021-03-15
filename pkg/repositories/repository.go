package repositories

import (
	"database/sql"
	"fmt"
	"github.com/malyg1n/sql-migrator/pkg/entities"
	"io/ioutil"
)

const (
	sqlScriptForMigrationTable = "sql_scripts/create-migrations-table.sql"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (repo *Repository) CheckOrCreateMigrationsTable() error {
	row, err := repo.db.Query(fmt.Sprintf("SELECT COUNT(id) from %s", migrationTableName))
	if err != nil || row == nil {
		data, err := ioutil.ReadFile(sqlScriptForMigrationTable)
		if err != nil {
			return err
		}
		_, err = repo.db.Exec(string(data))
		if err != nil {
			return err
		}
	}
	return nil
}

func (repo *Repository) GetMigrations() ([]*entities.MigrationEntity, error) {
	migrations := make([]*entities.MigrationEntity, 0)
	query := fmt.Sprintf("SELECT * from %s ORDER BY created_at DESC", migrationTableName)
	rows, err := repo.db.Query(query)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		me := &entities.MigrationEntity{}
		if err := rows.Scan(&me.Id, &me.Migration, &me.Version, &me.CreatedAt); err != nil {
			return nil, err
		}
		migrations = append(migrations, me)
	}

	return migrations, nil
}

func (repo *Repository) GetMigrationsByVersion(version uint) ([]*entities.MigrationEntity, error) {
	query := fmt.Sprintf("SELECT * FROM %s WHERE version=%d ORDER BY created_at DESC", migrationTableName, version)
	rows, err := repo.db.Query(query)
	if err != nil {
		return nil, err
	}
	migrations := make([]*entities.MigrationEntity, 0)
	for rows.Next() {
		me := &entities.MigrationEntity{}
		if err := rows.Scan(&me.Id, &me.Migration, &me.Version, &me.CreatedAt); err != nil {
			return nil, err
		}
		migrations = append(migrations, me)
	}

	return migrations, nil
}

func (repo *Repository) GetLatestVersionNumber() (uint, error) {
	var versionNumber uint
	me := &entities.MigrationEntity{}
	query := fmt.Sprintf("SELECT * FROM %s ORDER BY version DESC limit 1", migrationTableName)
	row := repo.db.QueryRow(query)
	if row == nil {
		versionNumber = 0
	} else {
		row.Scan(&me.Id, &me.Migration, &me.Version, &me.CreatedAt)
		if me.Version == 0 {
			return 0, nil
		}
		versionNumber = me.Version
	}
	return versionNumber, nil
}

func (repo *Repository) ApplyMigrationsUp(migrationName string, content string, version uint) error {
	migrationQuery := fmt.Sprintf("INSERT INTO %s (migration, version) VALUES ($1, $2)", migrationTableName)
	_, err := repo.db.Query(migrationQuery, migrationName, version)
	if err != nil {
		return err
	}
	_, err = repo.db.Query(content)
	if err != nil {
		return err
	}
	return nil
}

func (repo *Repository) ApplyMigrationsDown(migrationId uint, content string) error {
	_, err := repo.db.Query(content)
	if err != nil {
		return err
	}

	query := fmt.Sprintf("DELETE FROM %s WHERE id=%d", migrationTableName, migrationId)
	_, err = repo.db.Query(query)
	if err != nil {
		return err
	}

	return nil
}
