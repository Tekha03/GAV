package dog

import (
	"context"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) Create(ctx context.Context, dog *Dog) error {
	args := m.Called(ctx, dog)
	return args.Error(0)
}

func (m *MockRepository) Update(ctx context.Context, dog *Dog) error {
	args := m.Called(ctx, dog)
	return args.Error(0)
}

func (m *MockRepository) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockRepository) GetByOwnerID(ctx context.Context, ownerID uuid.UUID) ([]*Dog, error) {
	args := m.Called(ctx, ownerID)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).([]*Dog), args.Error(1)
}

func (m *MockRepository) GetByID(ctx context.Context, id uuid.UUID) (*Dog, error) {
	args := m.Called(ctx, id)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*Dog), args.Error(1)
}