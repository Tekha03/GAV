package dto

import "social_network/internal/auth"

type AuthResponse struct {
	Token *auth.AuthTokens `json:"token"`
}
