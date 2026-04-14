// internal/config/config.go
package config

import (
	"fmt"
	"net"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
    PostgresDSN          string        `yaml:"postgres_dsn"`
    RedisAddr            string        `yaml:"redis_addr"`
    SocialNetworkAddr    string        `yaml:"social_network_addr"`
    GRPCAddr             string        `yaml:"grpc_addr"`
    HTTPAddr             string        `yaml:"http_addr"`
    WSSAddr              string        `yaml:"wss_addr"`
    Env                  string        `yaml:"env"`
    LogLevel             string        `yaml:"log_level"`
    MigrationDir         string        `yaml:"migration_dir"`
    KafkaBrokers         []string      `yaml:"kafka_brokers"`
    KafkaTopic           string        `yaml:"kafka_topic"`
}

func Load() (*Config, error) {
    _ = godotenv.Load()

    cfg := &Config{
        Env:            getEnv("ENV", "development"),
        LogLevel:       getEnv("LOG_LEVEL", "info"),
        PostgresDSN:    getEnv("POSTGRES_DSN", ""),
        RedisAddr:      getEnv("REDIS_ADDR", "localhost:6379"),
        SocialNetworkAddr: getEnv("SOCIAL_NETWORK_ADDR", "localhost:9000"),
        GRPCAddr:       getEnv("GRPC_ADDR", ":9090"),
        HTTPAddr:       getEnv("HTTP_ADDR", ":8080"),
        WSSAddr:        getEnv("WSS_ADDR", ":8081"),
        MigrationDir:   getEnv("MIGRATION_DIR", "file://migrations"),
        KafkaBrokers:   strings.Split(getEnv("KAFKA_BROKERS", "localhost:9092"), ","),
        KafkaTopic:     getEnv("KAFKA_TOPIC", "events"),
    }

    if err := cfg.Validate(); err != nil {
        return nil, fmt.Errorf("config validation failed: %w", err)
    }

    return cfg, nil
}

func (c *Config) Validate() error {
    if c.PostgresDSN == "" {
        return fmt.Errorf("POSTGRES_DSN is required")
    }
    if _, err := parseDSN(c.PostgresDSN); err != nil {
        return fmt.Errorf("invalid POSTGRES_DSN: %w", err)
    }

    requiredAddrs := []string{c.RedisAddr, c.SocialNetworkAddr, c.GRPCAddr, c.HTTPAddr, c.WSSAddr}
    for _, addr := range requiredAddrs {
        if addr == "" {
            continue
        }
        if _, err := normalizeAddr(addr); err != nil {
            return fmt.Errorf("invalid address %s: %w", addr, err)
        }
    }

    validEnvs := map[string]bool{"development": true, "production": true, "test": true}
    if !validEnvs[c.Env] {
        return fmt.Errorf("invalid ENV: %s, must be 'development', 'production' or 'test'", c.Env)
    }

    validLogLevels := map[string]bool{"debug": true, "info": true, "warn": true, "error": true}
    if !validLogLevels[c.LogLevel] {
        return fmt.Errorf("invalid LOG_LEVEL: %s, must be 'debug', 'info', 'warn', 'error'", c.LogLevel)
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
        return nil, err
    }
    if u.Scheme != "postgres" && !strings.HasPrefix(dsn, "postgresql") {
        return nil, fmt.Errorf("unsupported DSN scheme")
    }
    return u, nil
}

func normalizeAddr(addr string) (string, error) {
    if !strings.Contains(addr, ":") {
        addr = ":" + addr
    }
    _, port, err := net.SplitHostPort(addr)
    if err != nil {
        return "", fmt.Errorf("invalid port")
    }
    if p, _ := strconv.Atoi(port); p < 1 || p > 65535 {
        return "", fmt.Errorf("port out of range")
    }
    return addr, nil
}
