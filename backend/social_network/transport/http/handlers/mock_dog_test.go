package handlers

import (
	"context"
	"social_network/internal/dog"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

type MockDogService struct {
	mock.Mock
}

func (m *MockDogService) Create(ctx context.Context, ownerID uuid.UUID, input dog.CreateDogInput) (*dog.Dog, error) {
	args := m.Called(ctx, ownerID, input)
	return args.Get(0).(*dog.Dog), args.Error(1)
}

func (m *MockDogService) Update(ctx context.Context, ownerID, dogID uuid.UUID, input dog.UpdateDogInput) error {
	args := m.Called(ctx, ownerID, dogID, input)
	return args.Error(0)
}

func (m *MockDogService) Delete(ctx context.Context, ownerID, dogID uuid.UUID) error {
	args := m.Called(ctx, ownerID, dogID)
	return args.Error(0)
}

func (m *MockDogService) UpdateLocation(ctx context.Context, ownerID, dogID uuid.UUID, lat, lon float64) error {
	args := m.Called(ctx, ownerID, dogID, lat, lon)
	return args.Error(0)
}

func (m *MockDogService) SetLocationVisibility(ctx context.Context, ownerID, dogID uuid.UUID, visible bool) error {
	args := m.Called(ctx, ownerID, dogID, visible)
	return args.Error(0)
}

func (m *MockDogService) GetPublic(ctx context.Context, dogID uuid.UUID) (*dog.Dog, error) {
	args := m.Called(ctx, dogID)
	return args.Get(0).(*dog.Dog), args.Error(1)
}

func (m *MockDogService) GetPrivate(ctx context.Context, ownerID, dogID uuid.UUID) (*dog.Dog, error) {
	args := m.Called(ctx, ownerID, dogID)
	return args.Get(0).(*dog.Dog), args.Error(1)
}
