package sqlite

import (
	"context"
	"testing"

	"social_network/internal/dog"
	"social_network/internal/user"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupUserRepository(t *testing.T) (*UserRepository, *gorm.DB) {
	t.Helper()

	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	require.NoError(t, err)
	require.NoError(t, db.AutoMigrate(&user.User{}, &dog.Dog{}))

	repo, err := NewUserRepository(db)
	require.NoError(t, err)

	return repo.(*UserRepository), db
}

func TestUserRepository_FindWalkingNearbyUsesOwnerLocation(t *testing.T) {
	repo, db := setupUserRepository(t)
	ctx := context.Background()

	lat := 55.751244
	lon := 37.618423
	ownerID := uuid.New()

	require.NoError(t, db.Create(&user.User{
		ID:             ownerID,
		Email:          "owner@gav.app",
		Password:       "hash",
		Role:           "user",
		Lat:            &lat,
		Lon:            &lon,
		LocationStatus: user.Walking,
		Visibility:     user.VisibilityEveryone,
	}).Error)

	d := &dog.Dog{
		ID:       uuid.New(),
		OwnerID:  ownerID,
		Name:     "Луна",
		Breed:    "Аусси",
		PhotoUrl: "/uploads/dogs/luna.jpg",
		Status:   dog.StatusFriendly,
		Age:      dog.AdultAge,
		Gender:   dog.Female,
	}
	require.NoError(t, db.Create(d).Error)

	dogs, err := repo.FindWalkingNearby(ctx, lat, lon, 1_000)

	require.NoError(t, err)
	require.Len(t, dogs, 1)
	require.Equal(t, d.ID, dogs[0].ID)
	require.NotNil(t, dogs[0].Lat)
	require.NotNil(t, dogs[0].Lon)
	require.Equal(t, lat, *dogs[0].Lat)
	require.Equal(t, lon, *dogs[0].Lon)
}
