package config

type DBConfig struct {
	PostgresDSN string
}

func loadDB() DBConfig {
	return DBConfig{
		PostgresDSN: getEnv("POSTGRES_DSN", ""),
	}
}
