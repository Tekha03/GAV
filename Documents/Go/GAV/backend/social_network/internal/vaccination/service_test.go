package vaccination

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestNewService(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		repo := new(MockVaccinationRepository)
		svc, err := NewService(repo)

		assert.NoError(t, err)
		assert.NotNil(t, svc)
	})

	t.Run("repo nil", func(t *testing.T) {
		svc, err := NewService(nil)

		assert.ErrorIs(t, err, ErrRepoNil)
		assert.Nil(t, svc)
	})
}

func TestService_Create(t *testing.T) {
	repo := new(MockVaccinationRepository)
	svc, _ := NewService(repo)

	ctx := context.Background()
	dogID := uuid.New()

	input := CreateVaccinationInput{
		Name:   "Rabies",
		DoneAt: time.Now(),
		Notes:  "test",
	}

	repo.On("Create", ctx, mock.AnythingOfType("*vaccination.Vaccination")).Return(nil)

	v, err := svc.Create(ctx, dogID, input)

	assert.NoError(t, err)
	assert.NotNil(t, v)
	assert.Equal(t, dogID, v.DogID)
	assert.Equal(t, input.Name, v.Name)

	repo.AssertExpectations(t)
}

func TestService_ListByDogID(t *testing.T) {
	repo := new(MockVaccinationRepository)
	svc, _ := NewService(repo)

	ctx := context.Background()
	dogID := uuid.New()

	expected := []*Vaccination{
		{ID: uuid.New(), DogID: dogID, Name: "Rabies"},
	}

	t.Run("success", func(t *testing.T) {
		repo.On("ListByDogID", ctx, dogID).Return(expected, nil)

		res, err := svc.ListByDogID(ctx, dogID)

		assert.NoError(t, err)
		assert.Equal(t, expected, res)
	})

	t.Run("empty dogID", func(t *testing.T) {
		res, err := svc.ListByDogID(ctx, uuid.Nil)

		assert.ErrorIs(t, err, ErrDogIDEmpty)
		assert.Nil(t, res)
	})
}

func TestService_Update(t *testing.T) {
	ctx := context.Background()
	vaccID := uuid.New()
	dogID := uuid.New()

	existing := &Vaccination{
		ID:    vaccID,
		DogID: dogID,
		Name:  "Old",
	}

	newName := "New"
	input := UpdateVaccinationInput{
		Name: &newName,
	}

	t.Run("success", func(t *testing.T) {
		repo := new(MockVaccinationRepository)
		svc, _ := NewService(repo)

		repo.On("GetByID", ctx, vaccID).Return(existing, nil)
		repo.On("Update", ctx, mock.AnythingOfType("*vaccination.Vaccination")).Return(nil)

		err := svc.Update(ctx, vaccID, dogID, input)

		assert.NoError(t, err)
		assert.Equal(t, newName, existing.Name)
	})

	t.Run("access denied", func(t *testing.T) {
		repo := new(MockVaccinationRepository)
		svc, _ := NewService(repo)

		otherDogID := uuid.New()

		repo.On("GetByID", ctx, vaccID).Return(existing, nil)

		err := svc.Update(ctx, vaccID, otherDogID, input)

		assert.ErrorIs(t, err, ErrVaccAccessDenied)
	})

	t.Run("repo error", func(t *testing.T) {
		repo := new(MockVaccinationRepository)
		svc, _ := NewService(repo)

		repo.On("GetByID", ctx, vaccID).Return((*Vaccination)(nil), ErrDBError)

		err := svc.Update(ctx, vaccID, dogID, input)

		assert.Error(t, err)
	})
}

func TestService_Delete(t *testing.T) {
	ctx := context.Background()
	vaccID := uuid.New()

	existing := &Vaccination{
		ID: vaccID,
	}

	t.Run("success", func(t *testing.T) {
		repo := new(MockVaccinationRepository)
		svc, _ := NewService(repo)

		repo.On("GetByID", ctx, vaccID).Return(existing, nil)
		repo.On("Delete", ctx, vaccID).Return(nil)

		err := svc.Delete(ctx, vaccID)

		assert.NoError(t, err)
	})

	t.Run("not found / repo error", func(t *testing.T) {
		repo := new(MockVaccinationRepository)
		svc, _ := NewService(repo)

		repo.On("GetByID", ctx, vaccID).Return((*Vaccination)(nil), ErrVaccinationNotFound)

		err := svc.Delete(ctx, vaccID)

		assert.Error(t, err)
	})

	t.Run("access denied", func(t *testing.T) {
		repo := new(MockVaccinationRepository)
		svc, _ := NewService(repo)

		wrong := &Vaccination{
			ID: uuid.New(),
		}

		repo.On("GetByID", ctx, vaccID).Return(wrong, nil)

		err := svc.Delete(ctx, vaccID)

		assert.ErrorIs(t, err, ErrVaccAccessDenied)
	})
}
