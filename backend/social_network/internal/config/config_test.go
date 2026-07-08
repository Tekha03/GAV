package config

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestErr(t *testing.T) {
	err := Err("test")

	assert.Error(t, err)
	assert.Equal(t, "config error: test", err.Error())
}

func TestLoadHTTP(t *testing.T) {
	t.Setenv("HTTP_PORT", "8080")
	cfg := loadHTTP()

	assert.Equal(t, "8080", cfg.Port)
}

func TestLoadDB(t *testing.T) {
	t.Setenv("DB_PATH", "./test.db")
	cfg := loadDB()

	assert.Equal(t, "./test.db", cfg.Path)
}

func TestLoadJWT(t *testing.T) {
	t.Setenv("JWT_SECRET", "secret")
	t.Setenv("JWT_TTL", "1h")

	cfg := loadJWT()

	assert.Equal(t, "secret", cfg.Secret)
	assert.Equal(t, time.Hour, cfg.TTL)
}

func TestLoadJWT_InvalidTTL(t *testing.T) {
	t.Setenv("JWT_SECRET", "secret")
	t.Setenv("JWT_TTL", "invalid")

	cfg := loadJWT()

	assert.Equal(t, "secret", cfg.Secret)
	assert.Equal(t, time.Duration(0), cfg.TTL)
}

func TestLoadStorage_Defaults(t *testing.T) {
	cfg := loadStorage()

	assert.Equal(t, "local", cfg.Type)
	assert.Equal(t, "./uploads", cfg.LocalPath)
}

func TestLoadStorage_FromEnv(t *testing.T) {
	t.Setenv("STORAGE_TYPE", "s3")
	t.Setenv("STORAGE_LOCAL_PATH", "/data")

	cfg := loadStorage()

	assert.Equal(t, "s3", cfg.Type)
	assert.Equal(t, "/data", cfg.LocalPath)
}

func TestGetEnv_Default(t *testing.T) {
	val := getEnv("NOT_EXISTING", "default")

	assert.Equal(t, "default", val)
}

func TestGetEnv_FromEnv(t *testing.T) {
	t.Setenv("TEST_KEY", "value")
	val := getEnv("TEST_KEY", "default")

	assert.Equal(t, "value", val)
}

func TestValidate_Success(t *testing.T) {
	cfg := &Config{
		HTTP: HTTPConfig{
			Port: "8080",
		},
		DB: DBConfig{
			Path: "./db.sqlite",
		},
		JWT: JWTConfig{
			Secret: "secret",
			TTL:    time.Hour,
		},
	}

	err := cfg.validate()
	assert.NoError(t, err)
}

func TestValidate_HTTPPortMissing(t *testing.T) {
	cfg := &Config{
		HTTP: HTTPConfig{},
		DB: DBConfig{
			Path: "./db",
		},
		JWT: JWTConfig{
			Secret: "secret",
			TTL:    time.Hour,
		},
	}

	err := cfg.validate()

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "HTTP_PORT")
}

func TestValidate_DBPathMissing(t *testing.T) {
	cfg := &Config{
		HTTP: HTTPConfig{
			Port: "8080",
		},
		DB: DBConfig{},
		JWT: JWTConfig{
			Secret: "secret",
			TTL:    time.Hour,
		},
	}

	err := cfg.validate()

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "DB_PATH")
}

func TestValidate_JWTSecretMissing(t *testing.T) {
	cfg := &Config{
		HTTP: HTTPConfig{
			Port: "8080",
		},
		DB: DBConfig{
			Path: "./db",
		},
		JWT: JWTConfig{
			TTL: time.Hour,
		},
	}

	err := cfg.validate()

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "JWT_SECRET")
}

func TestValidate_JWTTTLMissing(t *testing.T) {
	cfg := &Config{
		HTTP: HTTPConfig{
			Port: "8080",
		},
		DB: DBConfig{
			Path: "./db",
		},
		JWT: JWTConfig{
			Secret: "secret",
		},
	}

	err := cfg.validate()

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "JWT_TTL")
}

func TestLoad_Success(t *testing.T) {
	t.Setenv("HTTP_PORT", "8080")
	t.Setenv("DB_PATH", "./db.sqlite")
	t.Setenv("JWT_SECRET", "secret")
	t.Setenv("JWT_TTL", "1h")

	cfg, err := Load()

	assert.NoError(t, err)
	assert.Equal(t, "8080", cfg.HTTP.Port)
	assert.Equal(t, "./db.sqlite", cfg.DB.Path)
	assert.Equal(t, "secret", cfg.JWT.Secret)
	assert.Equal(t, time.Hour, cfg.JWT.TTL)
}

func TestLoad_Defaults(t *testing.T) {
	t.Setenv("HTTP_PORT", "")
	t.Setenv("DB_PATH", "")
	t.Setenv("JWT_SECRET", "")
	t.Setenv("JWT_TTL", "")

	cfg, err := Load()

	assert.NoError(t, err)
	assert.Equal(t, "8080", cfg.HTTP.Port)
	assert.Equal(t, "./dbserver/social.db", cfg.DB.Path)
	assert.Equal(t, "dev-secret-change-me", cfg.JWT.Secret)
	assert.Equal(t, 24*time.Hour, cfg.JWT.TTL)
}

func TestLoadEnvFrom(t *testing.T) {
	envFile := filepath.Join(t.TempDir(), ".env")
	os.WriteFile(envFile, []byte("TEST_ENV=hello"), 0644)

	err := loadEnvFromPath(envFile)

	assert.NoError(t, err)
	assert.Equal(t, "hello", os.Getenv("TEST_ENV"))
}
