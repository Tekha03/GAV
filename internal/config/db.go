package config

import "os"

type DBConfig struct {
	Path string
}

func loadDB() DBConfig {
	return DBConfig{
		Path: os.Getenv("DB_PATH"),
	}
}
