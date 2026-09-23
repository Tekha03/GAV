package user

import (
	"context"
	"time"

	"social_network/internal/dog"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) Create(ctx context.Context, user *User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockRepository) GetByEmail(ctx context.Context, email string) (*User, error) {
	args := m.Called(ctx, email)
	u := args.Get(0)
	if u == nil {
		return nil, args.Error(1)
	}
	return u.(*User), args.Error(1)
}

func (m *MockRepository) GetByID(ctx context.Context, id uuid.UUID) (*User, error) {
	args := m.Called(ctx, id)
	u := args.Get(0)
	if u == nil {
		return nil, args.Error(1)
	}
	return u.(*User), args.Error(1)
}

func (m *MockRepository) Update(ctx context.Context, user *User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockRepository) UpdateEmail(ctx context.Context, id uuid.UUID, email string) error {
	return m.Called(ctx, id, email).Error(0)
}

func (m *MockRepository) UpdatePassword(ctx context.Context, id uuid.UUID, oldHash, newHash string) error {
	return m.Called(ctx, id, oldHash, newHash).Error(0)
}

func (m *MockRepository) UpdateRole(ctx context.Context, id uuid.UUID, role string) error {
	return m.Called(ctx, id, role).Error(0)
}

func (m *MockRepository) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockRepository) FindWalkingNearby(ctx context.Context, centerLat, centerLon float64, radiusMeters float64, activeAfter time.Time) ([]*dog.Dog, error) {
	args := m.Called(ctx, centerLat, centerLon, radiusMeters, activeAfter)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).([]*dog.Dog), args.Error(1)
}

func (m *MockRepository) UpsertActiveWalkSession(ctx context.Context, session *WalkSession) error {
	args := m.Called(ctx, session)
	return args.Error(0)
}

func (m *MockRepository) EndActiveWalkSession(ctx context.Context, userID uuid.UUID, endedAt time.Time) error {
	args := m.Called(ctx, userID, endedAt)
	return args.Error(0)
}

func (m *MockRepository) UpdateActiveWalkVisibility(ctx context.Context, userID uuid.UUID, visibility LocationVisibility) error {
	args := m.Called(ctx, userID, visibility)
	return args.Error(0)
}
