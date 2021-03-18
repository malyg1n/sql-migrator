package migrations_store

import (
	"database/sql"
	"fmt"
	"github.com/malyg1n/sql-migrator/internal/entities"
)

type dBContract interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row
	Close() error
}

type MigrationsStore struct {
	db        dBContract
	tableName string
}

func NewStore(db dBContract, tableName string) *MigrationsStore {
	return &MigrationsStore{
		db:        db,
		tableName: tableName,
	}
}

func (s *MigrationsStore) CreateMigrationsTable(query string) error {
	_, err := s.db.Exec(query)
	return err
}

func (s *MigrationsStore) GetMigrations() ([]*entities.MigrationEntity, error) {
	var migrations []*entities.MigrationEntity
	query := fmt.Sprintf("SELECT * from %s ORDER BY created_at DESC", s.tableName)
	rows, err := s.db.Query(query)
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

func (s *MigrationsStore) GetMigrationsByVersion(version uint) ([]*entities.MigrationEntity, error) {
	query := fmt.Sprintf("SELECT * FROM %s WHERE version=%d ORDER BY created_at DESC", s.tableName, version)
	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}
	var migrations []*entities.MigrationEntity
	for rows.Next() {
		me := &entities.MigrationEntity{}
		if err := rows.Scan(&me.Id, &me.Migration, &me.Version, &me.CreatedAt); err != nil {
			return nil, err
		}
		migrations = append(migrations, me)
	}

	return migrations, nil
}

func (s *MigrationsStore) GetLatestVersionNumber() (uint, error) {
	var versionNumber uint
	me := &entities.MigrationEntity{}
	query := fmt.Sprintf("SELECT * FROM %s ORDER BY version DESC limit 1", s.tableName)
	row := s.db.QueryRow(query)
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

func (s *MigrationsStore) ApplyMigrationsUp(migrationName string, sqlQuery string, version uint) error {
	migrationQuery := fmt.Sprintf("INSERT INTO %s (migration, version) VALUES ($1, $2)", s.tableName)
	_, err := s.db.Exec(migrationQuery, migrationName, version)
	if err != nil {
		return err
	}
	_, err = s.db.Exec(sqlQuery)
	if err != nil {
		return err
	}
	return nil
}

func (s *MigrationsStore) ApplyMigrationsDown(migrationId uint, sqlQuery string) error {
	_, err := s.db.Exec(sqlQuery)
	if err != nil {
		return err
	}

	query := fmt.Sprintf("DELETE FROM %s WHERE id=%d", s.tableName, migrationId)
	_, err = s.db.Exec(query)
	if err != nil {
		return err
	}

	return nil
}
