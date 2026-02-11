package config

import (
	"os"
	"time"
)

type JWTConfig struct {
	Secret string
	TTL time.Duration
}

func loadJWT() JWTConfig {
	ttl, _ := time.ParseDuration(os.Getenv("JWT_TTL"))

	return JWTConfig{
		Secret: os.Getenv("JWT_SECRET"),
		TTL: ttl,
	}
}
