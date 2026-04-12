package handlers

import (
	"context"
	"social_network/internal/vaccination"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

type MockVaccinationService struct {
	mock.Mock
}

func (m *MockVaccinationService) Create(ctx context.Context, dogID uuid.UUID, input vaccination.CreateVaccinationInput) (*vaccination.Vaccination, error) {
	args := m.Called(ctx, dogID, input)
	return args.Get(0).(*vaccination.Vaccination), args.Error(1)
}

func (m *MockVaccinationService) ListByDogID(ctx context.Context, dogID uuid.UUID) ([]*vaccination.Vaccination, error) {
	args := m.Called(ctx, dogID)
	return args.Get(0).([]*vaccination.Vaccination), args.Error(1)
}

func (m *MockVaccinationService) Update(ctx context.Context, vaccinationID, dogID uuid.UUID, input vaccination.UpdateVaccinationInput) error {
	args := m.Called(ctx, vaccinationID, dogID, input)
	return args.Error(0)
}

func (m *MockVaccinationService) Delete(ctx context.Context, vaccinationID uuid.UUID) error {
	args := m.Called(ctx, vaccinationID)
	return args.Error(0)
}
