// internal/config/config.go
package config

import (
	"net"
	"net/url"
	"os"
	apperrors "shared/app_errors"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	PostgresDSN       string   `yaml:"postgres_dsn"`
	RedisAddr         string   `yaml:"redis_addr"`
	SocialNetworkAddr string   `yaml:"social_network_addr"`
	GRPCAddr          string   `yaml:"grpc_addr"`
	HTTPAddr          string   `yaml:"http_addr"`
	WSSAddr           string   `yaml:"wss_addr"`
	Env               string   `yaml:"env"`
	LogLevel          string   `yaml:"log_level"`
	MigrationDir      string   `yaml:"migration_dir"`
	KafkaEnabled      bool     `yaml:"kafka_enabled"`
	KafkaBrokers      []string `yaml:"kafka_brokers"`
	KafkaTopic        string   `yaml:"kafka_topic"`
	JWTSecret         string   `yaml:"jwt_secret"`
}

func Load() (*Config, error) {
	_ = godotenv.Load()

	cfg := &Config{
		Env:               getEnv("ENV", "development"),
		LogLevel:          getEnv("LOG_LEVEL", "info"),
		PostgresDSN:       getEnv("POSTGRES_DSN", ""),
		RedisAddr:         getEnv("REDIS_ADDR", "localhost:6379"),
		SocialNetworkAddr: getEnv("SOCIAL_NETWORK_ADDR", "localhost:9000"),
		GRPCAddr:          getEnv("GRPC_ADDR", ":9090"),
		HTTPAddr:          getEnv("HTTP_ADDR", ":8082"),
		WSSAddr:           getEnv("WSS_ADDR", ":8081"),
		MigrationDir:      getEnv("MIGRATION_DIR", "file://migrations"),
		KafkaEnabled:      getEnv("KAFKA_ENABLED", "false") == "true",
		KafkaBrokers:      strings.Split(getEnv("KAFKA_BROKERS", "localhost:9092"), ","),
		KafkaTopic:        getEnv("KAFKA_TOPIC", "events"),
		JWTSecret:         getEnv("JWT_SECRET", ""),
	}

	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

func (c *Config) Validate() error {
	if c.PostgresDSN == "" {
		return apperrors.New(apperrors.Validation, "POSTGRES_DSN is required", apperrors.WithDetail("field", "POSTGRES_DSN"))
	}
	if c.JWTSecret == "" {
		return apperrors.New(apperrors.Validation, "JWT_SECRET is required", apperrors.WithDetail("field", "JWT_SECRET"))
	}
	if _, err := parseDSN(c.PostgresDSN); err != nil {
		return apperrors.Wrap(apperrors.Validation, "invalid POSTGRES_DSN", err, apperrors.WithDetail("field", "POSTGRES_DSN"))
	}

	requiredAddrs := []string{c.RedisAddr, c.SocialNetworkAddr, c.GRPCAddr, c.HTTPAddr, c.WSSAddr}
	for _, addr := range requiredAddrs {
		if addr == "" {
			continue
		}
		if _, err := normalizeAddr(addr); err != nil {
			return apperrors.Wrap(apperrors.Validation, "invalid service address", err, apperrors.WithDetail("address", addr))
		}
	}

	validEnvs := map[string]bool{"development": true, "production": true, "test": true}
	if !validEnvs[c.Env] {
		return apperrors.New(apperrors.Validation, "invalid ENV", apperrors.WithDetail("value", c.Env))
	}

	validLogLevels := map[string]bool{"debug": true, "info": true, "warn": true, "error": true}
	if !validLogLevels[c.LogLevel] {
		return apperrors.New(apperrors.Validation, "invalid LOG_LEVEL", apperrors.WithDetail("value", c.LogLevel))
	}

	return nil
}

func getEnv(key, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return strings.TrimSpace(val)
	}
	return defaultVal
}

func parseDSN(dsn string) (*url.URL, error) {
	u, err := url.Parse(dsn)
	if err != nil {
		return nil, apperrors.Wrap(apperrors.Validation, "invalid DSN", err)
	}
	if u.Scheme != "postgres" && !strings.HasPrefix(dsn, "postgresql") {
		return nil, apperrors.New(apperrors.Validation, "unsupported DSN scheme")
	}
	return u, nil
}

func normalizeAddr(addr string) (string, error) {
	if !strings.Contains(addr, ":") {
		addr = ":" + addr
	}
	_, port, err := net.SplitHostPort(addr)
	if err != nil {
		return "", apperrors.Wrap(apperrors.Validation, "invalid port", err)
	}
	if p, _ := strconv.Atoi(port); p < 1 || p > 65535 {
		return "", apperrors.New(apperrors.Validation, "port out of range", apperrors.WithDetail("port", port))
	}
	return addr, nil
}
