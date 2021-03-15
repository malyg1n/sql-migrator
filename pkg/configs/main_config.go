package configs

import "os"

type MainConfig struct {
	MigrationsPath string
}

func NewMainConfig() *MainConfig {
	return &MainConfig{
		MigrationsPath: os.Getenv("MIGRATIONS_PATH"),
	}
}
