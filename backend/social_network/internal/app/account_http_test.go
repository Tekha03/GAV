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
	"social_network/internal/testdb"
	"social_network/internal/user"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

func TestAccountUpdateHTTPPostgres(t *testing.T) {
	t.Setenv("ENV", "production")
	const secret = "account-test-secret"
	cfg := &config.Config{
		HTTP: config.HTTPConfig{Port: "8080"}, GRPC: config.GRPCConfig{Addr: ":9000"},
		DB:      config.DBConfig{PostgresDSN: testdb.DSN(t)},
		JWT:     config.JWTConfig{Secret: secret, TTL: time.Hour},
		Storage: config.StorageConfig{Type: "local", LocalPath: t.TempDir()},
	}
	application, err := NewApp(context.Background(), cfg)
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, application.Shutdown(context.Background())) })
	hash, err := bcrypt.GenerateFromPassword([]byte("old-password"), bcrypt.MinCost)
	require.NoError(t, err)
	owner, err := application.Services.User.Create(context.Background(), "owner@test.local", string(hash))
	require.NoError(t, err)
	other, err := application.Services.User.Create(context.Background(), "other@test.local", string(hash))
	require.NoError(t, err)
	admin, err := application.Services.User.Create(context.Background(), "admin@test.local", string(hash))
	require.NoError(t, err)
	require.NoError(t, application.sqlDB.Model(&user.User{}).Where("id = ?", admin.ID).Update("role", "admin").Error)
	lat := 55.75
	require.NoError(t, application.sqlDB.Model(&user.User{}).Where("id = ?", owner.ID).Updates(map[string]any{"lat": lat, "visibility": user.VisibilityFollowersOnly}).Error)

	token := func(id uuid.UUID) string {
		t.Helper()
		value, err := auth.GenerateAccessToken(id, "user", auth.JWTConfig{Secret: []byte(secret), TTL: time.Hour})
		require.NoError(t, err)
		return value
	}
	request := func(method, path, bearer, body string) int {
		t.Helper()
		req := httptest.NewRequest(method, path, bytes.NewBufferString(body))
		if bearer != "" {
			req.Header.Set("Authorization", "Bearer "+bearer)
		}
		w := httptest.NewRecorder()
		application.Server.Handler.ServeHTTP(w, req)
		return w.Code
	}
	path := "/api/v1/users/" + owner.ID.String()
	require.Equal(t, http.StatusUnauthorized, request(http.MethodPut, path, "", `{"email":"new@test.local"}`))
	require.Equal(t, http.StatusForbidden, request(http.MethodPut, path, token(other.ID), `{"email":"new@test.local"}`))
	require.Equal(t, http.StatusBadRequest, request(http.MethodPut, path, token(owner.ID), `{}`))
	require.Equal(t, http.StatusBadRequest, request(http.MethodPut, path, token(owner.ID), `{"email":"bad"}`))
	require.Equal(t, http.StatusBadRequest, request(http.MethodPut, path, token(owner.ID), `{"role":"admin"}`))
	require.Equal(t, http.StatusBadRequest, request(http.MethodPut, path, token(owner.ID), `{"password":"injected"}`))
	require.Equal(t, http.StatusConflict, request(http.MethodPut, path, token(owner.ID), `{"email":"other@test.local"}`))
	missingID := uuid.New()
	require.Equal(t, http.StatusNotFound, request(http.MethodPut, "/api/v1/users/"+missingID.String(), token(missingID), `{"email":"new@test.local"}`))
	require.Equal(t, http.StatusNoContent, request(http.MethodPut, path, token(owner.ID), `{"email":"new@test.local"}`))

	var account user.User
	require.NoError(t, application.sqlDB.First(&account, "id = ?", owner.ID).Error)
	require.Equal(t, "new@test.local", account.Email)
	require.Equal(t, string(hash), account.Password)
	require.Equal(t, "user", account.Role)
	require.NotNil(t, account.Lat)
	require.Equal(t, lat, *account.Lat)
	require.Equal(t, user.VisibilityFollowersOnly, account.Visibility)

	require.Equal(t, http.StatusForbidden, request(http.MethodPut, path+"/role", token(other.ID), `{"role":"admin"}`))
	require.Equal(t, http.StatusNoContent, request(http.MethodPut, path+"/role", token(admin.ID), `{"role":"admin"}`))
	require.NoError(t, application.sqlDB.First(&account, "id = ?", owner.ID).Error)
	require.Equal(t, "admin", account.Role)

	require.Equal(t, http.StatusBadRequest, request(http.MethodPut, "/api/v1/users/me/password", token(owner.ID), `{"current_password":"old-password","new_password":"short"}`))
	require.Equal(t, http.StatusUnauthorized, request(http.MethodPut, "/api/v1/users/me/password", token(owner.ID), `{"current_password":"wrong-password","new_password":"new-password"}`))
	require.Equal(t, http.StatusNoContent, request(http.MethodPut, "/api/v1/users/me/password", token(owner.ID), `{"current_password":"old-password","new_password":"new-password"}`))
	require.NoError(t, application.sqlDB.First(&account, "id = ?", owner.ID).Error)
	require.NoError(t, bcrypt.CompareHashAndPassword([]byte(account.Password), []byte("new-password")))
	require.Equal(t, "new@test.local", account.Email)
}
