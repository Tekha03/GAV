package auth

import "time"

type JWTConfig struct {
	Secret []byte
	TTL time.Duration
}
