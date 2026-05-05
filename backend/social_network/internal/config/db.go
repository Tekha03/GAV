package config

import "os"

type DBConfig struct {
	DSN string
}

func loadDB() DBConfig {
	return DBConfig{
		DSN: os.Getenv("POSTGRES_DSN"),
	}
}
