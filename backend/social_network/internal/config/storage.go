package config

import "os"

type StorageConfig struct {
	Type		string
	LocalPath	string
}

func loadStorage() StorageConfig {
	return StorageConfig{
		Type: getEnv("STORAGE_TYPE", "local"),
		LocalPath: getEnv("STORAGE_LOCAL_PATH", "./uploads"),
	}
}

func getEnv(key, defaultValue string) string {
	if env := os.Getenv(key); env != "" {
		return env
	}
	return defaultValue
}
