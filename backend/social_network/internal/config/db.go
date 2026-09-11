package config

type DBConfig struct {
	Driver      string
	Path        string
	PostgresDSN string
}

func loadDB() DBConfig {
	return DBConfig{
		Driver:      getEnv("DB_DRIVER", "sqlite"),
		Path:        getEnv("DB_PATH", "./dbserver/social.db"),
		PostgresDSN: getEnv("POSTGRES_DSN", ""),
	}
}
