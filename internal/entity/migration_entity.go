package entity

import "time"

// MigrationEntity for database
type MigrationEntity struct {
	ID        uint      `db:"id"`
	Migration string    `db:"migration"`
	Version   uint      `db:"version"`
	CreatedAt time.Time `db:"created_at"`
	Query     string
}

// NewMigrationEntity returns the new instance
func NewMigrationEntity(migration, query string, version uint) *MigrationEntity {
	return &MigrationEntity{
		Migration: migration,
		Query:     query,
		Version:   version,
	}
}
