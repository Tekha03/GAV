package app

import "errors"

var (
	ErrConfigNil        = errors.New("app: config is nil")
	ErrDBPathEmpty      = errors.New("app: db path is empty")
	ErrDBDriverEmpty    = errors.New("app: db driver is empty")
	ErrPostgresDSNEmpty = errors.New("app: postgres dsn is empty")
	ErrJWTSecretEmpty   = errors.New("app: jwt secret is empty")
)
