package app

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"social_network/internal/auth"
	"social_network/internal/config"
	"social_network/internal/device"
	"social_network/internal/testdb"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestDeviceTokenRegistrationHTTPPostgres(t *testing.T) {
	t.Setenv("ENV", "production")
	const secret = "device-http-test-secret"
	cfg := &config.Config{
		HTTP: config.HTTPConfig{Port: "8080"}, GRPC: config.GRPCConfig{Addr: ":9000"},
		DB:      config.DBConfig{PostgresDSN: testdb.DSN(t)},
		JWT:     config.JWTConfig{Secret: secret, TTL: time.Hour},
		Storage: config.StorageConfig{Type: "local", LocalPath: t.TempDir()},
	}
	application, err := NewApp(context.Background(), cfg)
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, application.Shutdown(context.Background())) })
	first, err := application.Services.User.Create(context.Background(), uuid.NewString()+"@test.local", "hash")
	require.NoError(t, err)
	second, err := application.Services.User.Create(context.Background(), uuid.NewString()+"@test.local", "hash")
	require.NoError(t, err)
	request := func(method string, userID uuid.UUID, body string) int {
		t.Helper()
		req := httptest.NewRequest(method, "/api/v1/devices/tokens", bytes.NewBufferString(body))
		if userID != uuid.Nil {
			token, err := auth.GenerateAccessToken(userID, "user", auth.JWTConfig{Secret: []byte(secret), TTL: time.Hour})
			require.NoError(t, err)
			req.Header.Set("Authorization", "Bearer "+token)
		}
		w := httptest.NewRecorder()
		application.Server.Handler.ServeHTTP(w, req)
		return w.Code
	}
	require.Equal(t, http.StatusUnauthorized, request(http.MethodPost, uuid.Nil, `{"token":"fcm-token"}`))
	require.Equal(t, http.StatusBadRequest, request(http.MethodPost, first.ID, `{"token":""}`))
	require.Equal(t, http.StatusNoContent, request(http.MethodPost, first.ID, `{"token":"fcm-token"}`))
	require.Equal(t, http.StatusNoContent, request(http.MethodPost, second.ID, `{"token":"fcm-token"}`))
	var firstTokens []device.DeviceToken
	require.NoError(t, application.sqlDB.Where("user_id = ?", first.ID).Find(&firstTokens).Error)
	require.Empty(t, firstTokens)
	var secondTokens []device.DeviceToken
	require.NoError(t, application.sqlDB.Where("user_id = ?", second.ID).Find(&secondTokens).Error)
	require.Len(t, secondTokens, 1)
	require.Equal(t, http.StatusNoContent, request(http.MethodDelete, first.ID, `{"token":"fcm-token"}`))
	require.Equal(t, http.StatusNoContent, request(http.MethodDelete, second.ID, `{"token":"fcm-token"}`))
	secondTokens = nil
	require.NoError(t, application.sqlDB.Where("user_id = ?", second.ID).Find(&secondTokens).Error)
	require.Empty(t, secondTokens)
}
