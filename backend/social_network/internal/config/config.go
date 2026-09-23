package config

type Config struct {
	HTTP     HTTPConfig
	GRPC     GRPCConfig
	DB       DBConfig
	JWT      JWTConfig
	Storage  StorageConfig
	Firebase FirebaseConfig
}

func Load() (*Config, error) {
	loadEnv()

	cfg := &Config{
		HTTP:     loadHTTP(),
		GRPC:     loadGRPC(),
		DB:       loadDB(),
		JWT:      loadJWT(),
		Storage:  loadStorage(),
		Firebase: loadFirebase(),
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
	if c.GRPC.Addr == "" {
		return Err("GRPC_ADDR is required")
	}
	if c.DB.PostgresDSN == "" {
		return Err("POSTGRES_DSN is required")
	}
	if c.JWT.Secret == "" {
		return Err("JWT_SECRET is required")
	}
	if c.JWT.TTL == 0 {
		return Err("JWT_TTL is required")
	}
	if c.Firebase.Enabled && c.Firebase.CredentialsFile == "" {
		return Err("FIREBASE_CREDENTIALS_FILE is required when FIREBASE_ENABLED=true")
	}
	return nil
}
