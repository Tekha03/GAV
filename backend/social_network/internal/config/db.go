package config

type DBConfig struct {
	Path string
}

func loadDB() DBConfig {
	return DBConfig{
		Path: getEnv("DB_PATH", "./dbserver/social.db"),
	}
}
