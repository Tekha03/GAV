package grpc

import "errors"

var (
	ErrIvalidAuthorization = errors.New("invalid authorization header")
	ErrInvalidSignMethod   = errors.New("unexpected signing method")
	ErrInvalidToken        = errors.New("invalid token")
	ErrInvalidSubject      = errors.New("invalid subject")
	ErrMissingMetadata     = errors.New("missing metadata")
	ErrMissingAuthMetadata = errors.New("missing authorization metadata")
)
