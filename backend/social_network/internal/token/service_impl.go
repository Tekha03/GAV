package token

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"time"

	"github.com/google/uuid"
)

type service struct {
	repo Repository
}

func NewService(repo Repository) (TokenService, error) {
	if repo == nil {
		return nil, ErrRepoNil
	}

	return &service{repo: repo}, nil
}

func (s *service) CreateRefresh(ctx context.Context, userID uuid.UUID) (string, error) {
	plain, err := generateRandomRefresh()
	if err != nil {
		return "", err
	}

	token := &RefreshToken{
		ID:        uuid.New(),
		UserID:    userID,
		TokenHash: refreshHash(plain),
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour),
		CreatedAt: time.Now(),
	}

	if err := s.repo.Create(ctx, token); err != nil {
		return "", err
	}

	return plain, nil
}

func (s *service) ValidateAndRotate(ctx context.Context, refresh string) (uuid.UUID, string, error) {
	hash := refreshHash(refresh)
	refreshToken, err := s.repo.GetByHash(ctx, hash)
	if err != nil || refreshToken == nil || refreshToken.Revoked || time.Now().After(refreshToken.ExpiresAt) {
		return uuid.Nil, "", ErrInvalidRefresh
	}

	if err := s.repo.Revoke(ctx, hash); err != nil {
		return uuid.Nil, "", err
	}
	newPlain, err := s.CreateRefresh(ctx, refreshToken.UserID)
	if err != nil {
		return uuid.Nil, "", err
	}

	return refreshToken.UserID, newPlain, nil
}

func (s *service) Revoke(ctx context.Context, refresh string) error {
	return s.repo.Revoke(ctx, refreshHash(refresh))
}

func refreshHash(refresh string) string {
	sum := sha256.Sum256([]byte(refresh))
	return hex.EncodeToString(sum[:])
}

func (s *service) RevokeAllForUser(ctx context.Context, userID uuid.UUID) error {
	return s.repo.RevokeAllForUser(ctx, userID)
}

func generateRandomRefresh() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}

	return base64.URLEncoding.EncodeToString(b), nil

}
