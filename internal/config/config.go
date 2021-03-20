package config

import "os"

// Config is configuration for app
type Config struct {
	DbDriver            string
	DbConnectionsString string
	MigrationsPath      string
}

// NewConfig returns new Config instance
func NewConfig() *Config {
	return &Config{
		DbDriver:            os.Getenv("DB_DRIVER"),
		DbConnectionsString: os.Getenv("DB_DSN"),
		MigrationsPath:      os.Getenv("MIGRATIONS_PATH"),
	}
}
