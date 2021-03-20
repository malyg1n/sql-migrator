package config

import "os"

type Config struct {
	DbDriver            string
	DbConnectionsString string
	MigrationsPath      string
	PrepareScriptsPath  string
}

func NewConfig() *Config {
	return &Config{
		DbDriver:            os.Getenv("DB_DRIVER"),
		DbConnectionsString: os.Getenv("DB_DSN"),
		MigrationsPath:      os.Getenv("MIGRATIONS_PATH"),
		PrepareScriptsPath:  "prepare",
	}
}
