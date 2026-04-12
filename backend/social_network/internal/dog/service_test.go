package dog

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type testEnv struct {
	service DogService
	repo	*MockRepository
}

func setup(t *testing.T) *testEnv {
	repo := &MockRepository{}

	service, err := NewService(repo)
	require.NoError(t, err)

	return &testEnv{
		service: service,
		repo: repo,
	}
}

func TestNewService(t *testing.T) {
	env := setup(t)
	assert.NotNil(t, env.service)
}

func TestNewService_RepoNil(t *testing.T) {
	service, err := NewService(nil)

	assert.Error(t, err)
	assert.Nil(t, service)
	assert.Equal(t, ErrRepoNil, err)
}

func TestCreate_Success(t *testing.T) {
	env := setup(t)

	ctx := context.Background()
	ownerID := uuid.New()

	input := CreateDogInput{
		Name: "Buddy",
		Breed: "Labrador",
		Gender: Male,
		Status: StatusFriendly,
		Age: AdultAge,
		PhotoUrl: "url",
	}

	env.repo.
		On("Create", ctx, mock.AnythingOfType("*dog.Dog")).
		Return(nil)

	dog, err := env.service.Create(ctx, ownerID, input)

	assert.NoError(t, err)
	assert.Equal(t, ownerID, dog.OwnerID)
}

func TestUpdate_Success(t *testing.T) {
	env := setup(t)

	ctx := context.Background()
	ownerID := uuid.New()
	dogID := uuid.New()

	name := "NewName"

	dog := &Dog{
		ID: dogID,
		OwnerID: ownerID,
		Name: "Old",
	}

	env.repo.
		On("GetByID", ctx, dogID).
		Return(dog, nil)

	env.repo.
		On("Update", ctx, dog).
		Return(nil)

	input := UpdateDogInput{
		Name: &name,
	}

	err := env.service.Update(ctx, ownerID, dogID, input)

	assert.NoError(t, err)
	assert.Equal(t, "NewName", dog.Name)
}

func TestUpdate_AccessDenied(t *testing.T) {
	env := setup(t)
	ctx := context.Background()

	dogID := uuid.New()
	dog := &Dog{
		ID: dogID,
		OwnerID: uuid.New(),
	}

	env.repo.
		On("GetByID", ctx, dogID).
		Return(dog, nil)

	err := env.service.Update(ctx, uuid.New(), dogID, UpdateDogInput{})

	assert.Error(t, err)
	assert.Equal(t, ErrDogAccessDenied, err)
}

func TestDelete_Success(t *testing.T) {
	env := setup(t)

	ctx := context.Background()
	ownerID := uuid.New()
	dogID := uuid.New()

	dog := &Dog{
		ID: dogID,
		OwnerID: ownerID,
	}

	env.repo.On("GetByID", ctx, dogID).Return(dog, nil)
	env.repo.On("Delete", ctx, dogID).Return(nil)

	err := env.service.Delete(ctx, ownerID, dogID)
	assert.NoError(t, err)
}

func TestGetPrivate_AccessDenied(t *testing.T) {
	env := setup(t)
	ctx := context.Background()

	dogID := uuid.New()
	ownerID := uuid.New()

	dog := &Dog{
		ID: dogID,
		OwnerID: ownerID,
	}

	env.repo.On("GetByID", ctx, dogID).Return(dog, nil)
	_, err := env.service.GetPrivate(ctx, dogID, uuid.New())

	assert.Error(t, err)
	assert.Equal(t, ErrDogAccessDenied, err)
}
