package device

import (
	"context"
	"strings"

	"github.com/google/uuid"
)

type Service struct{ repo Repository }

func NewService(repo Repository) *Service { return &Service{repo: repo} }

func (s *Service) Register(ctx context.Context, userID uuid.UUID, token string) error {
	token = strings.TrimSpace(token)
	if userID == uuid.Nil || len(token) == 0 || len(token) > 4096 {
		return ErrInvalidToken
	}
	return s.repo.SaveForUser(ctx, userID, token)
}

func (s *Service) Unregister(ctx context.Context, userID uuid.UUID, token string) error {
	token = strings.TrimSpace(token)
	if userID == uuid.Nil || len(token) == 0 || len(token) > 4096 {
		return ErrInvalidToken
	}
	return s.repo.DeleteForUser(ctx, userID, token)
}
