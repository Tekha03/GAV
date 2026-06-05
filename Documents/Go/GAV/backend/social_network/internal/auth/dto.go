package auth

import "github.com/google/uuid"

type UserInfo struct {
	ID    uuid.UUID
	Email string
}

type AuthTokens struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}
