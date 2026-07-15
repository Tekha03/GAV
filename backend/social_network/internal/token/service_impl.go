package token

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
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

	hash, err := bcrypt.GenerateFromPassword([]byte(plain), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	token := &RefreshToken{
		ID:        uuid.New(),
		UserID:    userID,
		TokenHash: string(hash),
		ExpiresAt: time.Now().Add(7*24 + time.Hour),
		CreatedAt: time.Now(),
	}

	if err := s.repo.Create(ctx, token); err != nil {
		return "", err
	}

	return plain, nil
}

func (s *service) ValidateAndRotate(ctx context.Context, refresh string) (uuid.UUID, string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(refresh), bcrypt.DefaultCost)
	if err != nil {
		return uuid.Nil, "", err
	}

	refreshToken, err := s.repo.GetByHash(ctx, string(hash))
	if err != nil || refreshToken == nil || refreshToken.Revoked || time.Now().After(refreshToken.ExpiresAt) {
		return uuid.Nil, "", ErrInvalidRefresh
	}

	s.repo.Revoke(ctx, string(hash))
	newPlain, err := s.CreateRefresh(ctx, refreshToken.UserID)
	if err != nil {
		return uuid.Nil, "", err
	}

	return refreshToken.UserID, newPlain, err
}

func (s *service) Revoke(ctx context.Context, refresh string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(refresh), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	return s.repo.Revoke(ctx, string(hash))
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
