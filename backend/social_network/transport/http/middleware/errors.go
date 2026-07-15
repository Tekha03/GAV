package middleware

import "errors"

var (
	ErrUnauthorized = errors.New("middleware: unauthorized")
	ErrForbidden    = errors.New("forbidden: insufficient permissions")
	ErrNotOwner     = errors.New("forbidden: not owner")
	ErrInvalidID    = errors.New("middleware: invalid id")
)
