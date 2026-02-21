package config

type Config struct {
	HTTP HTTPConfig
	DB DBConfig
	JWT JWTConfig
}

func Load() (*Config, error) {
	loadEnv()

	cfg := &Config{
		HTTP: loadHTTP(),
		DB: loadDB(),
		JWT: loadJWT(),
	}

	if err := cfg.validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

func (c *Config) validate() error {
	if c.HTTP.Port == "" {
		return Err("HTTP_PORT is required")
	}
	if c.DB.Path == "" {
		return Err("DB_PATH is required")
	}
	if c.JWT.Secret == "" {
		return Err("JWT_SECRET is required")
	}
	if c.JWT.TTL == 0 {
		return Err("JWT_TTL is required")
	}
	return nil
}
