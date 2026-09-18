package app

import (
	"context"
	"testing"
	"time"

	"social_network/internal/config"
	"social_network/internal/testdb"

	"github.com/stretchr/testify/require"
)

func TestFirebaseEnabledWithMissingCredentialsFailsStartup(t *testing.T) {
	t.Setenv("ENV", "production")
	cfg := &config.Config{
		HTTP: config.HTTPConfig{Port: "8080"}, GRPC: config.GRPCConfig{Addr: ":9000"},
		DB:       config.DBConfig{PostgresDSN: testdb.DSN(t)},
		JWT:      config.JWTConfig{Secret: "test-secret", TTL: time.Hour},
		Storage:  config.StorageConfig{Type: "local", LocalPath: t.TempDir()},
		Firebase: config.FirebaseConfig{Enabled: true, CredentialsFile: t.TempDir() + "/missing.json"},
	}
	application, err := NewApp(context.Background(), cfg)
	require.Nil(t, application)
	require.ErrorContains(t, err, "firebase credentials file unavailable")
}
