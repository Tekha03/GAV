package config

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFirebaseConfigRequiresCredentialsWhenEnabled(t *testing.T) {
	t.Setenv("FIREBASE_ENABLED", "true")
	t.Setenv("FIREBASE_CREDENTIALS_FILE", "")
	cfg := loadFirebase()
	require.True(t, cfg.Enabled)
	require.Equal(t, "", cfg.CredentialsFile)
	base := &Config{HTTP: HTTPConfig{Port: "8080"}, GRPC: GRPCConfig{Addr: ":9000"}, DB: DBConfig{PostgresDSN: "postgres://localhost/test"}, JWT: JWTConfig{Secret: "test", TTL: 1}}
	base.Firebase = cfg
	require.ErrorContains(t, base.validate(), "FIREBASE_CREDENTIALS_FILE")
	base.Firebase.CredentialsFile = "/tmp/service-account.json"
	require.NoError(t, base.validate())
}
