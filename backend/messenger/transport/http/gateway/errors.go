package gateway

import "errors"

var (
	errMissingToken            = errors.New("missing authorization token")
	errInvalidToken            = errors.New("invalid authorization token")
	errUnexpectedSigningMethod = errors.New("unexpected token signing method")
	errUnauthorized            = errors.New("unauthorized")
	errForbidden               = errors.New("forbidden")
)
