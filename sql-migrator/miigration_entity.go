package sql_migrator

import "time"

type MigrationEntity struct {
	Id        uint      `db:"id"`
	Migration string    `db:"migration"`
	Version   uint      `db:"version"`
	CreatedAt time.Time `db:"created_at"`
}
