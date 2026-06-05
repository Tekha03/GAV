package handlers

import (
	"context"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"

	"social_network/internal/comment"
)

type MockCommentService struct {
	mock.Mock
}

func (m *MockCommentService) Create(ctx context.Context, userID, postID uuid.UUID, content string) error {
	args := m.Called(ctx, userID, postID, content)
	return args.Error(0)
}

func (m *MockCommentService) GetByID(ctx context.Context, id uuid.UUID) (*comment.Comment, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*comment.Comment), args.Error(1)
}

func (m *MockCommentService) ListByPostID(ctx context.Context, postID uuid.UUID) ([]comment.Comment, error) {
	args := m.Called(ctx, postID)
	return args.Get(0).([]comment.Comment), args.Error(1)
}

func (m *MockCommentService) Delete(ctx context.Context, userID, commentID uuid.UUID) error {
	args := m.Called(ctx, userID, commentID)
	return args.Error(0)
}
