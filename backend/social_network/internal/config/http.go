package config

type HTTPConfig struct {
	Port string
}

func loadHTTP() HTTPConfig {
	return HTTPConfig{
		Port: getEnv("HTTP_PORT", "8080"),
	}
}
