package post

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) Create(ctx context.Context, post *Post) error {
	args := m.Called(ctx, post)
	return args.Error(0)
}

func (m *MockRepository) GetByID(ctx context.Context, postID uuid.UUID) (*Post, error) {
	args := m.Called(ctx, postID)
	post := args.Get(0)

	if post == nil {
		return nil, args.Error(1)
	}

	return post.(*Post), args.Error(1)
}

func (m *MockRepository) ListByUser(ctx context.Context, authorID uuid.UUID) ([]*Post, error) {
	args := m.Called(ctx, authorID)
	post_list := args.Get(0)

	if post_list == nil {
		return nil, args.Error(1)
	}

	return post_list.([]*Post), args.Error(1)
}

func (m *MockRepository) ListFeed(ctx context.Context, userID uuid.UUID, before time.Time, limit int) ([]*Post, error) {
	args := m.Called(ctx, userID, before, limit)
	post_list := args.Get(0)

	if post_list == nil {
		return nil, args.Error(1)
	}

	return post_list.([]*Post), args.Error(1)
}

func (m *MockRepository) Delete(ctx context.Context, postID uuid.UUID) error {
	args := m.Called(ctx, postID)
	return args.Error(0)
}
