package comment

import (
	"context"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) Create(ctx context.Context, comment *Comment) error {
	args := m.Called(ctx, comment)
	return args.Error(0)
}

func (m *MockRepository) GetByID(ctx context.Context, id uuid.UUID) (*Comment, error) {
	args := m.Called(ctx, id)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*Comment), args.Error(1)
}

func (m *MockRepository) ListByPostID(ctx context.Context, postID uuid.UUID) ([]Comment, error) {
	args := m.Called(ctx, postID)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).([]Comment), args.Error(1)
}

func (m *MockRepository) Delete(ctx context.Context, userID, commentID uuid.UUID) error {
	args := m.Called(ctx, userID, commentID)
	return args.Error(0)
}
