package config

import "os"

type HTTPConfig struct {
	Port string
}

func loadHTTP() HTTPConfig {
	return HTTPConfig{
		Port: os.Getenv("HTTP_POST"),
	}
}
