package config

import "errors"

func Err(msg string) error {
	return errors.New("config error: " + msg)
}
