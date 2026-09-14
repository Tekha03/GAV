package client

import (
	authv1 "api/gen/auth/v1"
	socialv1 "api/gen/social/v1"
	"context"

	"github.com/google/uuid"
)

type NotifiClient interface {
	SendNewMessage(ctx context.Context, toUserID uuid.UUID, senderName, body, chatID string) error
	Close() error
}

type SocialClient interface {
	GetUserProfile(ctx context.Context, userID uuid.UUID) (*socialv1.UserProfile, error)
	Login(ctx context.Context, email, password string) (*authv1.LoginResponse, error)
	Close() error
}
