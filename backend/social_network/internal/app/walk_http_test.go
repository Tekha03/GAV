package app

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"social_network/internal/auth"
	"social_network/internal/config"
	"social_network/internal/follow"
	"social_network/internal/testdb"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestWalkHTTPPrivacyFlow(t *testing.T) {
	t.Setenv("ENV", "production")
	const secret = "walk-http-test-secret"
	cfg := &config.Config{
		HTTP: config.HTTPConfig{Port: "8080"}, GRPC: config.GRPCConfig{Addr: ":9000"},
		DB:      config.DBConfig{PostgresDSN: testdb.DSN(t)},
		JWT:     config.JWTConfig{Secret: secret, TTL: time.Hour},
		Storage: config.StorageConfig{Type: "local", LocalPath: t.TempDir()},
	}
	application, err := NewApp(context.Background(), cfg)
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, application.Shutdown(context.Background())) })
	owner, err := application.Services.User.Create(context.Background(), uuid.NewString()+"@test.local", "hash")
	require.NoError(t, err)
	follower, err := application.Services.User.Create(context.Background(), uuid.NewString()+"@test.local", "hash")
	require.NoError(t, err)
	outsider, err := application.Services.User.Create(context.Background(), uuid.NewString()+"@test.local", "hash")
	require.NoError(t, err)
	token := func(id uuid.UUID) string {
		t.Helper()
		value, err := auth.GenerateAccessToken(id, "user", auth.JWTConfig{Secret: []byte(secret), TTL: time.Hour})
		require.NoError(t, err)
		return value
	}
	request := func(method, path, bearer string, body []byte) *httptest.ResponseRecorder {
		t.Helper()
		req := httptest.NewRequest(method, path, bytes.NewReader(body))
		if bearer != "" {
			req.Header.Set("Authorization", "Bearer "+bearer)
		}
		result := httptest.NewRecorder()
		application.Server.Handler.ServeHTTP(result, req)
		return result
	}
	path := "/api/v1/walks/nearby?lat=55.75&lon=37.61&radius=1000"
	start := request(http.MethodPost, "/api/v1/walks/start", token(owner.ID), []byte(`{"latitude":55.75,"longitude":37.61,"visibility":1}`))
	require.Equal(t, http.StatusCreated, start.Code, start.Body.String())
	require.Equal(t, http.StatusConflict, request(http.MethodPost, "/api/v1/walks/start", token(owner.ID), []byte(`{"latitude":55.75,"longitude":37.61,"visibility":1}`)).Code)
	require.Equal(t, http.StatusUnauthorized, request(http.MethodGet, path, "", nil).Code)
	require.Equal(t, http.StatusForbidden, request(http.MethodPut, "/api/v1/users/"+owner.ID.String()+"/location-visibility", token(outsider.ID), []byte(`{"visibility":2}`)).Code)
	require.Equal(t, http.StatusOK, request(http.MethodGet, "/api/v1/walks/current", token(owner.ID), nil).Code)

	decodeResults := func(token string) []map[string]any {
		t.Helper()
		result := request(http.MethodGet, path, token, nil)
		require.Equal(t, http.StatusOK, result.Code, result.Body.String())
		var items []map[string]any
		require.NoError(t, json.Unmarshal(result.Body.Bytes(), &items))
		return items
	}
	require.Empty(t, decodeResults(token(follower.ID)))
	require.Empty(t, decodeResults(token(outsider.ID)))
	require.NoError(t, application.sqlDB.Create(&follow.Follow{FollowerID: follower.ID, FollowingID: owner.ID}).Error)
	require.Len(t, decodeResults(token(follower.ID)), 1)
	require.Empty(t, decodeResults(token(outsider.ID)))
	require.Equal(t, http.StatusNoContent, request(http.MethodPut, "/api/v1/users/"+owner.ID.String()+"/location-visibility", token(owner.ID), []byte(`{"visibility":2}`)).Code)
	require.Empty(t, decodeResults(token(follower.ID)))
	require.Equal(t, http.StatusOK, request(http.MethodGet, "/api/v1/walks/current", token(owner.ID), nil).Code)
	require.Equal(t, http.StatusOK, request(http.MethodPatch, "/api/v1/walks/current/visibility", token(owner.ID), []byte(`{"visibility":1}`)).Code)
	require.Len(t, decodeResults(token(follower.ID)), 1)

	require.Equal(t, http.StatusOK, request(http.MethodPatch, "/api/v1/walks/current/visibility", token(owner.ID), []byte(`{"visibility":2}`)).Code)
	require.Empty(t, decodeResults(token(follower.ID)))
	require.Equal(t, http.StatusOK, request(http.MethodGet, "/api/v1/walks/current", token(owner.ID), nil).Code)
	require.Equal(t, http.StatusNoContent, request(http.MethodPost, "/api/v1/walks/current/stop", token(owner.ID), nil).Code)
	require.Equal(t, http.StatusNotFound, request(http.MethodGet, "/api/v1/walks/current", token(owner.ID), nil).Code)
}
