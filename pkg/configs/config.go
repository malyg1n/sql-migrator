package configs

type Config struct {
	DB   *DBConfig
	Main *MainConfig
}

func NewConfig() *Config {
	return &Config{
		DB:   NewDBConfig(),
		Main: NewMainConfig(),
	}
}
