package app

import "errors"

var (
	ErrConfigNil 	 	 = errors.New("app: config is nil")
	ErrDBDSNEmpty 		= errors.New("app: db dsn is empty")
	ErrJWTSecretEmpty 	= errors.New("app: jwt secret is empty")
)
