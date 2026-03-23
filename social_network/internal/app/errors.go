package app

import "errors"

var (
	ErrConfigNil 	 	 = errors.New("app: config is nil")
	ErrDBPathEmpty 		= errors.New("app: db path is empty")
	ErrJWTSecretEmpty 	= errors.New("app: jwt secret is empty")
)
