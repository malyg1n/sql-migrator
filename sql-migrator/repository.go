package sql_migrator

import (
	"fmt"
)

type Repository struct {
	db DBContract
}

func NewRepository(db DBContract) *Repository {
	return &Repository{
		db: db,
	}
}

func (repo *Repository) CreateMigrationsTable(query string) error {
	_, err := repo.db.Exec(query)
	return err
}

func (repo *Repository) GetMigrations() ([]*MigrationEntity, error) {
	migrations := make([]*MigrationEntity, 0)
	query := fmt.Sprintf("SELECT * from %s ORDER BY created_at DESC", migrationTableName)
	rows, err := repo.db.Query(query)
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

func (repo *Repository) GetMigrationsByVersion(version uint) ([]*MigrationEntity, error) {
	query := fmt.Sprintf("SELECT * FROM %s WHERE version=%d ORDER BY created_at DESC", migrationTableName, version)
	rows, err := repo.db.Query(query)
	if err != nil {
		return nil, err
	}
	migrations := make([]*MigrationEntity, 0)
	for rows.Next() {
		me := &MigrationEntity{}
		if err := rows.Scan(&me.Id, &me.Migration, &me.Version, &me.CreatedAt); err != nil {
			return nil, err
		}
		migrations = append(migrations, me)
	}

	return migrations, nil
}

func (repo *Repository) GetLatestVersionNumber() (uint, error) {
	var versionNumber uint
	me := &MigrationEntity{}
	query := fmt.Sprintf("SELECT * FROM %s ORDER BY version DESC limit 1", migrationTableName)
	row := repo.db.QueryRow(query)
	if row == nil {
		versionNumber = 0
	} else {
		err := row.Scan(&me.Id, &me.Migration, &me.Version, &me.CreatedAt)
		if me.Version == 0 || err != nil {
			return 0, nil
		}
		versionNumber = me.Version
	}
	return versionNumber, nil
}

func (repo *Repository) ApplyMigrationsUp(migrationName string, content string, version uint) error {
	migrationQuery := fmt.Sprintf("INSERT INTO %s (migration, version) VALUES ($1, $2)", migrationTableName)
	_, err := repo.db.Exec(migrationQuery, migrationName, version)
	if err != nil {
		return err
	}
	_, err = repo.db.Exec(content)
	if err != nil {
		return err
	}
	return nil
}

func (repo *Repository) ApplyMigrationsDown(migrationId uint, content string) error {
	_, err := repo.db.Exec(content)
	if err != nil {
		return err
	}

	query := fmt.Sprintf("DELETE FROM %s WHERE id=%d", migrationTableName, migrationId)
	_, err = repo.db.Exec(query)
	if err != nil {
		return err
	}

	return nil
}
