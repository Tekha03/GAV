package dbserver

import "errors"

var (
	ErrCannotOpen = errors.New("cannot open postgres")
)
