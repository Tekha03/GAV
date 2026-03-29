package vaccination

import (
	"context"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

type MockVaccinationRepository struct {
	mock.Mock
}

func (m *MockVaccinationRepository) Create(ctx context.Context, v *Vaccination) error {
	args := m.Called(ctx, v)
	return args.Error(0)
}

func (m *MockVaccinationRepository) Update(ctx context.Context, v *Vaccination) error {
	args := m.Called(ctx, v)
	return args.Error(0)
}

func (m *MockVaccinationRepository) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockVaccinationRepository) ListByDogID(ctx context.Context, dogID uuid.UUID) ([]*Vaccination, error) {
	args := m.Called(ctx, dogID)
	return args.Get(0).([]*Vaccination), args.Error(1)
}

func (m *MockVaccinationRepository) GetByID(ctx context.Context, id uuid.UUID) (*Vaccination, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*Vaccination), args.Error(1)
}
