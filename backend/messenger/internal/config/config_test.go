package config

import (
	"os"
	"strings"
	"testing"
)

func TestGetEnv(t *testing.T) {
	key := "MESSENGER_TEST_ENV"
	_ = os.Unsetenv(key)
	if got := getEnv(key, "default"); got != "default" {
		t.Fatalf("unexpected default value: %s", got)
	}
	if err := os.Setenv(key, "  value  "); err != nil {
		t.Fatalf("Setenv returned error: %v", err)
	}
	if got := getEnv(key, "default"); got != "value" {
		t.Fatalf("unexpected trimmed value: %s", got)
	}
}

func TestParseDSN(t *testing.T) {
	if _, err := parseDSN("postgres://user:pass@localhost:5432/db"); err != nil {
		t.Fatalf("parseDSN returned error: %v", err)
	}
	if _, err := parseDSN("mysql://user:pass@localhost:3306/db"); err == nil {
		t.Fatal("expected error for unsupported scheme")
	}
}

func TestNormalizeAddr(t *testing.T) {
	addr, err := normalizeAddr("9090")
	if err != nil {
		t.Fatalf("normalizeAddr returned error: %v", err)
	}
	if addr != ":9090" {
		t.Fatalf("unexpected normalized addr: %s", addr)
	}
	addr, err = normalizeAddr("localhost:8080")
	if err != nil {
		t.Fatalf("normalizeAddr returned error: %v", err)
	}
	if addr != "localhost:8080" {
		t.Fatalf("unexpected normalized addr: %s", addr)
	}
	if _, err := normalizeAddr("localhost:not-a-port"); err == nil {
		t.Fatal("expected error for invalid port")
	}
	if _, err := normalizeAddr("localhost:70000"); err == nil {
		t.Fatal("expected error for out of range port")
	}
}

func TestValidate(t *testing.T) {
	cfg := &Config{
		PostgresDSN:       "postgres://user:pass@localhost:5432/db",
		RedisAddr:         "localhost:6379",
		SocialNetworkAddr: "localhost:9000",
		GRPCAddr:          ":9090",
		HTTPAddr:          ":8080",
		WSSAddr:           ":8081",
		Env:               "development",
		LogLevel:          "info",
	}
	if err := cfg.Validate(); err != nil {
		t.Fatalf("Validate returned error: %v", err)
	}

	cases := []Config{
		{PostgresDSN: "", RedisAddr: "localhost:6379", SocialNetworkAddr: "localhost:9000", GRPCAddr: ":9090", HTTPAddr: ":8080", WSSAddr: ":8081", Env: "development", LogLevel: "info"},
		{PostgresDSN: "mysql://bad", RedisAddr: "localhost:6379", SocialNetworkAddr: "localhost:9000", GRPCAddr: ":9090", HTTPAddr: ":8080", WSSAddr: ":8081", Env: "development", LogLevel: "info"},
		{PostgresDSN: "postgres://ok", RedisAddr: "bad-port", SocialNetworkAddr: "localhost:9000", GRPCAddr: ":9090", HTTPAddr: ":8080", WSSAddr: ":8081", Env: "development", LogLevel: "info"},
		{PostgresDSN: "postgres://ok", RedisAddr: "localhost:6379", SocialNetworkAddr: "localhost:9000", GRPCAddr: ":9090", HTTPAddr: ":8080", WSSAddr: ":8081", Env: "qa", LogLevel: "info"},
		{PostgresDSN: "postgres://ok", RedisAddr: "localhost:6379", SocialNetworkAddr: "localhost:9000", GRPCAddr: ":9090", HTTPAddr: ":8080", WSSAddr: ":8081", Env: "development", LogLevel: "trace"},
	}

	for i := range cases {
		if err := cases[i].Validate(); err == nil {
			t.Fatalf("case %d should fail", i)
		}
	}
}

func TestLoad(t *testing.T) {
	pairs := map[string]string{
		"ENV":                 "test",
		"LOG_LEVEL":           "debug",
		"POSTGRES_DSN":        "postgres://user:pass@localhost:5432/db",
		"REDIS_ADDR":          "localhost:6379",
		"SOCIAL_NETWORK_ADDR": "localhost:9000",
		"GRPC_ADDR":           ":9090",
		"HTTP_ADDR":           ":8080",
		"WSS_ADDR":            ":8081",
		"MIGRATION_DIR":       "file://migrations",
		"KAFKA_BROKERS":       "localhost:9092,localhost:9093",
		"KAFKA_TOPIC":         "events",
	}
	for key, value := range pairs {
		if err := os.Setenv(key, value); err != nil {
			t.Fatalf("Setenv %s returned error: %v", key, err)
		}
	}

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load returned error: %v", err)
	}
	if cfg.Env != "test" || cfg.LogLevel != "debug" {
		t.Fatalf("unexpected config values: %+v", cfg)
	}
	if len(cfg.KafkaBrokers) != 2 || strings.TrimSpace(cfg.KafkaBrokers[1]) != "localhost:9093" {
		t.Fatalf("unexpected KafkaBrokers: %+v", cfg.KafkaBrokers)
	}
}
