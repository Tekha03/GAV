package config

import (
	"testing"
)

func validConfig() *Config {
	return &Config{
		PostgresDSN:       "postgres://user:password@localhost:5432/messenger",
		RedisAddr:         "localhost:6379",
		SocialNetworkAddr: "localhost:9000",
		GRPCAddr:          ":9090",
		HTTPAddr:          ":8082",
		WSSAddr:           ":8081",
		Env:               "test",
		LogLevel:          "info",
		JWTSecret:         "secret",
	}
}

func TestConfigValidate(t *testing.T) {
	if err := validConfig().Validate(); err != nil {
		t.Fatalf("Validate() error = %v", err)
	}
	for _, test := range []struct {
		name   string
		mutate func(*Config)
	}{
		{name: "missing postgres", mutate: func(c *Config) { c.PostgresDSN = "" }},
		{name: "missing jwt", mutate: func(c *Config) { c.JWTSecret = "" }},
		{name: "unsupported dsn", mutate: func(c *Config) { c.PostgresDSN = "mysql://localhost/db" }},
		{name: "bad address", mutate: func(c *Config) { c.HTTPAddr = "host:not-a-port" }},
		{name: "bad environment", mutate: func(c *Config) { c.Env = "demo" }},
		{name: "bad log level", mutate: func(c *Config) { c.LogLevel = "verbose" }},
	} {
		t.Run(test.name, func(t *testing.T) {
			config := validConfig()
			test.mutate(config)
			if err := config.Validate(); err == nil {
				t.Fatal("Validate() error = nil")
			}
		})
	}
}

func TestNormalizeAddrAndParseDSN(t *testing.T) {
	for input, want := range map[string]string{
		"8080":           ":8080",
		":8080":          ":8080",
		"localhost:8080": "localhost:8080",
	} {
		got, err := normalizeAddr(input)
		if err != nil || got != want {
			t.Fatalf("normalizeAddr(%q) = %q, %v", input, got, err)
		}
	}
	for _, input := range []string{"0", "65536", "localhost"} {
		if _, err := normalizeAddr(input); err == nil {
			t.Fatalf("normalizeAddr(%q) error = nil", input)
		}
	}
	if _, err := parseDSN("postgresql://localhost/db"); err != nil {
		t.Fatalf("parseDSN(postgresql) error = %v", err)
	}
}

func TestLoadReadsEnvironmentAndDefaults(t *testing.T) {
	t.Setenv("POSTGRES_DSN", "postgres://user:password@localhost:5432/messenger")
	t.Setenv("JWT_SECRET", "secret")
	t.Setenv("ENV", "test")
	t.Setenv("LOG_LEVEL", "debug")
	t.Setenv("KAFKA_ENABLED", "true")
	t.Setenv("KAFKA_BROKERS", "first:9092,second:9092")

	config, err := Load()
	if err != nil {
		t.Fatal(err)
	}
	if !config.KafkaEnabled || len(config.KafkaBrokers) != 2 || config.LogLevel != "debug" {
		t.Fatalf("Load() = %+v", config)
	}
}
