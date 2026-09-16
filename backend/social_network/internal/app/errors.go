package app

import "errors"

var (
	ErrConfigNil        = errors.New("app: config is nil")
	ErrPostgresDSNEmpty = errors.New("app: postgres dsn is empty")
	ErrJWTSecretEmpty   = errors.New("app: jwt secret is empty")
)
