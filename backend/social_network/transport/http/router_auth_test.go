package http

import (
	"bytes"
	"context"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"social_network/internal/auth"
	"social_network/transport/http/handlers"
)

type routerAuthService struct{}

func (routerAuthService) Register(context.Context, string, string) (*auth.AuthTokens, error) {
	return &auth.AuthTokens{AccessToken: "registered", RefreshToken: "refresh"}, nil
}

func (routerAuthService) Login(_ context.Context, _ string, password string) (*auth.AuthTokens, error) {
	if password != "correct-password" {
		return nil, auth.ErrInvalidCredentials
	}
	return &auth.AuthTokens{AccessToken: "logged-in", RefreshToken: "refresh"}, nil
}

func (routerAuthService) Me(context.Context, uuid.UUID) (*auth.UserInfo, error) {
	return nil, nil
}

func (routerAuthService) Logout(context.Context, string) error { return nil }

func (routerAuthService) Refresh(context.Context, string) (*auth.AuthTokens, error) {
	return &auth.AuthTokens{AccessToken: "refreshed", RefreshToken: "refresh"}, nil
}

func TestAuthRoutesWithoutAccessToken(t *testing.T) {
	authHandler, err := handlers.NewAuthHandler(routerAuthService{})
	require.NoError(t, err)

	requireAccessToken := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Header.Get("Authorization") != "Bearer valid" {
				http.Error(w, "missing access token", http.StatusUnauthorized)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
	router := NewRouter(
		Handlers{Auth: authHandler},
		RouterDeps{AuthMW: requireAccessToken},
		slog.New(slog.NewTextHandler(io.Discard, nil)),
	)

	request := func(method, path, body string) *httptest.ResponseRecorder {
		t.Helper()
		req := httptest.NewRequest(method, path, bytes.NewBufferString(body))
		response := httptest.NewRecorder()
		router.ServeHTTP(response, req)
		return response
	}

	registered := request(http.MethodPost, "/api/v1/auth/register", `{"email":"new@example.com","password":"correct-password"}`)
	require.Equal(t, http.StatusCreated, registered.Code, registered.Body.String())
	require.Contains(t, registered.Body.String(), "registered")

	loggedIn := request(http.MethodPost, "/api/v1/auth/login", `{"email":"new@example.com","password":"correct-password"}`)
	require.Equal(t, http.StatusOK, loggedIn.Code, loggedIn.Body.String())
	require.Contains(t, loggedIn.Body.String(), "logged-in")

	badLogin := request(http.MethodPost, "/api/v1/auth/login", `{"email":"new@example.com","password":"wrong-password"}`)
	require.Equal(t, http.StatusUnauthorized, badLogin.Code, badLogin.Body.String())
	require.Contains(t, badLogin.Body.String(), "AUTH_CREDENTIALS_INVALID")

	refreshed := request(http.MethodPost, "/api/v1/auth/refresh", `{"refresh_token":"refresh"}`)
	require.Equal(t, http.StatusOK, refreshed.Code, refreshed.Body.String())
	require.Contains(t, refreshed.Body.String(), "refreshed")

	protected := request(http.MethodGet, "/api/v1/auth/me", "")
	require.Equal(t, http.StatusUnauthorized, protected.Code)
}
