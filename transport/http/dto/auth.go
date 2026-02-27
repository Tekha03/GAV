package dto

import "gav/internal/auth"

type AuthResponse struct {
	Token *auth.AuthTokens `json:"token"`
}
