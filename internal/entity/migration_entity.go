package entity

import "time"

type MigrationEntity struct {
	Id        uint      `db:"id"`
	Migration string    `db:"migration"`
	Version   uint      `db:"version"`
	CreatedAt time.Time `db:"created_at"`
	Query     string
}

func NewMigrationEntity(migration, query string, version uint) *MigrationEntity {
	return &MigrationEntity{
		Migration: migration,
		Query:     query,
		Version:   version,
	}
}
