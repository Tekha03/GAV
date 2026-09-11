package dbserver

import "errors"

var (
	ErrUnsupportedDriver = errors.New("unsupported database driver")
)
