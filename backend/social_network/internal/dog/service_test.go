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

func TestUpdateLocation(t *testing.T) {
	env := setup(t)
	ctx := context.Background()

	ownerID := uuid.New()
	dogID := uuid.New()

	dog := &Dog{
		ID: dogID,
		OwnerID: ownerID,
	}

	env.repo.On("GetByID", ctx, dogID).Return(dog, nil)
	env.repo.On("Update", ctx, dog).Return(nil)

	err := env.service.UpdateLocation(ctx, ownerID, dogID, UpdateLocationInput{Latitude: 10, Longitude: 20})

	assert.NoError(t, err)
	assert.NotNil(t, dog.Lat)
	assert.NotNil(t, dog.Lon)
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

func TestFindDogsNearby_Success(t *testing.T) {
	env := setup(t)

	ctx := context.Background()
	userID := uuid.New()
	otherUserID := uuid.New()

	dogsFromRepo := []*Dog{
		{
			ID:              uuid.New(),
			OwnerID:         otherUserID,
			Visibility: 	 1,
		},
	}

	env.repo.
		On("FindWalkingNearby", ctx, 0.0, 0.0, 1000.0).
		Return(dogsFromRepo, nil)

	dogs, err := env.service.FindDogsNearby(ctx, userID, 0.0, 0.0, 1000.0)

	require.NoError(t, err)
	require.Len(t, dogs, 1)
	assert.Equal(t, otherUserID, dogs[0].OwnerID)
}

func TestFindDogsNearby_ExcludeOwnDogs(t *testing.T) {
	env := setup(t)
	ctx := context.Background()
	userID := uuid.New()

	dogsFromRepo := []*Dog{
		{
			ID:              uuid.New(),
			OwnerID:         userID,
			Visibility: 	 1,
		},
	}

	env.repo.
		On("FindWalkingNearby", ctx, 0.0, 0.0, 1000.0).
		Return(dogsFromRepo, nil)

	dogs, err := env.service.FindDogsNearby(ctx, userID, 0, 0, 1000)

	require.NoError(t, err)
	require.Len(t, dogs, 0)
}

func TestFindDogsNearby_FilterInvisible(t *testing.T) {
	env := setup(t)

	ctx := context.Background()
	userID := uuid.New()

	dogsFromRepo := []*Dog{
		{
			ID:              uuid.New(),
			OwnerID:         uuid.New(),
			Visibility: 	 0,
		},
	}

	env.repo.
		On("FindWalkingNearby", ctx, 0.0, 0.0, 1000.0).
		Return(dogsFromRepo, nil)

	dogs, err := env.service.FindDogsNearby(ctx, userID, 0, 0, 1000)

	require.NoError(t, err)
	require.Len(t, dogs, 0)
}

func TestFindDogsNearby_RepoError(t *testing.T) {
	env := setup(t)

	ctx := context.Background()
	userID := uuid.New()

	env.repo.
		On("FindWalkingNearby", ctx, 0.0, 0.0, 1000.0).
		Return(nil, assert.AnError)

	dogs, err := env.service.FindDogsNearby(ctx, userID, 0, 0, 1000)

	require.Error(t, err)
	assert.Nil(t, dogs)
}
