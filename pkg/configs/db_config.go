package configs

import "os"

// Config for database
type DBConfig struct {
	// for postgres and mysql
	Driver   string
	Host     string
	Port     string
	Name     string
	User     string
	Password string

	// Only for postgres
	SSLMode string

	// Only for sqlite
	File  string
	Cache string
	Mode  string

	// Another drivers
	DSN string

	// Query for create
	QueryForCreateTable string
}

// Create new DBConfig
func NewDBConfig() *DBConfig {
	return &DBConfig{
		Driver:              os.Getenv("DB_DRIVER"),
		Host:                os.Getenv("DB_HOST"),
		Port:                os.Getenv("DB_PORT"),
		Name:                os.Getenv("DB_NAME"),
		User:                os.Getenv("DB_USER"),
		Password:            os.Getenv("DB_PASSWORD"),
		SSLMode:             os.Getenv("DB_SSL_MODE"),
		File:                os.Getenv("DB_FILE"),
		DSN:                 os.Getenv("DB_DSN"),
		QueryForCreateTable: getSqlQueryToCreateMigrationsTable(),
	}
}

func getSqlQueryToCreateMigrationsTable() string {
	return ""
}
