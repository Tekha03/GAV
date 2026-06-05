package config

import "time"

type JWTConfig struct {
	Secret string
	TTL    time.Duration
}

func loadJWT() JWTConfig {
	ttl, _ := time.ParseDuration(getEnv("JWT_TTL", "24h"))

	return JWTConfig{
		Secret: getEnv("JWT_SECRET", "dev-secret-change-me"),
		TTL:    ttl,
	}
}
