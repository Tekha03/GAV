package sqlite

import (
	"context"
	"testing"
	"time"

	"social_network/internal/dog"
	"social_network/internal/user"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupUserRepository(t *testing.T) (*UserRepository, *gorm.DB) {
	t.Helper()

	db, err := gorm.Open(sqlite.Open("file:"+uuid.NewString()+"?mode=memory&cache=shared"), &gorm.Config{})
	require.NoError(t, err)
	require.NoError(t, db.AutoMigrate(&user.User{}, &user.WalkSession{}, &dog.Dog{}))

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
	now := time.Now()

	require.NoError(t, db.Create(&user.User{
		ID:         ownerID,
		Email:      "owner@gav.app",
		Password:   "hash",
		Role:       "user",
		Visibility: user.VisibilityEveryone,
	}).Error)
	require.NoError(t, db.Create(&user.WalkSession{
		ID:         uuid.New(),
		UserID:     ownerID,
		Lat:        lat,
		Lon:        lon,
		Visibility: user.VisibilityEveryone,
		StartedAt:  now,
		UpdatedAt:  now,
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

	dogs, err := repo.FindWalkingNearby(ctx, lat, lon, 1_000, time.Now().Add(-2*time.Hour))

	require.NoError(t, err)
	require.Len(t, dogs, 1)
	require.Equal(t, d.ID, dogs[0].ID)
	require.NotNil(t, dogs[0].Lat)
	require.NotNil(t, dogs[0].Lon)
	require.Equal(t, lat, *dogs[0].Lat)
	require.Equal(t, lon, *dogs[0].Lon)
}

func TestUserRepository_FindWalkingNearbyFiltersByRadiusAndFreshness(t *testing.T) {
	repo, db := setupUserRepository(t)
	ctx := context.Background()

	centerLat := 55.751244
	centerLon := 37.618423
	fresh := time.Now()
	stale := time.Now().Add(-3 * time.Hour)

	nearOwnerID := uuid.New()
	nearLat := centerLat + 0.001
	nearLon := centerLon
	require.NoError(t, db.Create(&user.User{
		ID:         nearOwnerID,
		Email:      "near@gav.app",
		Password:   "hash",
		Role:       "user",
		Visibility: user.VisibilityEveryone,
	}).Error)
	require.NoError(t, db.Create(&user.WalkSession{
		ID:         uuid.New(),
		UserID:     nearOwnerID,
		Lat:        nearLat,
		Lon:        nearLon,
		Visibility: user.VisibilityEveryone,
		StartedAt:  fresh,
		UpdatedAt:  fresh,
	}).Error)
	require.NoError(t, db.Create(&dog.Dog{
		ID:       uuid.New(),
		OwnerID:  nearOwnerID,
		Name:     "Рядом",
		Breed:    "Корги",
		PhotoUrl: "/uploads/dogs/near.jpg",
		Status:   dog.StatusFriendly,
		Age:      dog.AdultAge,
		Gender:   dog.Female,
	}).Error)

	farOwnerID := uuid.New()
	farLat := centerLat + 0.05
	farLon := centerLon
	require.NoError(t, db.Create(&user.User{
		ID:         farOwnerID,
		Email:      "far@gav.app",
		Password:   "hash",
		Role:       "user",
		Visibility: user.VisibilityEveryone,
	}).Error)
	require.NoError(t, db.Create(&user.WalkSession{
		ID:         uuid.New(),
		UserID:     farOwnerID,
		Lat:        farLat,
		Lon:        farLon,
		Visibility: user.VisibilityEveryone,
		StartedAt:  fresh,
		UpdatedAt:  fresh,
	}).Error)
	require.NoError(t, db.Create(&dog.Dog{
		ID:       uuid.New(),
		OwnerID:  farOwnerID,
		Name:     "Далеко",
		Breed:    "Хаски",
		PhotoUrl: "/uploads/dogs/far.jpg",
		Status:   dog.StatusFriendly,
		Age:      dog.AdultAge,
		Gender:   dog.Male,
	}).Error)

	staleOwnerID := uuid.New()
	staleLat := centerLat
	staleLon := centerLon
	require.NoError(t, db.Create(&user.User{
		ID:         staleOwnerID,
		Email:      "stale@gav.app",
		Password:   "hash",
		Role:       "user",
		Visibility: user.VisibilityEveryone,
	}).Error)
	require.NoError(t, db.Create(&user.WalkSession{
		ID:         uuid.New(),
		UserID:     staleOwnerID,
		Lat:        staleLat,
		Lon:        staleLon,
		Visibility: user.VisibilityEveryone,
		StartedAt:  stale,
		UpdatedAt:  stale,
	}).Error)
	require.NoError(t, db.Create(&dog.Dog{
		ID:       uuid.New(),
		OwnerID:  staleOwnerID,
		Name:     "Старая",
		Breed:    "Шпиц",
		PhotoUrl: "/uploads/dogs/stale.jpg",
		Status:   dog.StatusFriendly,
		Age:      dog.AdultAge,
		Gender:   dog.Female,
	}).Error)

	dogs, err := repo.FindWalkingNearby(ctx, centerLat, centerLon, 500, time.Now().Add(-2*time.Hour))

	require.NoError(t, err)
	require.Len(t, dogs, 1)
	require.Equal(t, "Рядом", dogs[0].Name)
	require.NotNil(t, dogs[0].DistanceMeters)
	require.Less(t, *dogs[0].DistanceMeters, 500.0)
}

func TestUserRepository_FindWalkingNearbySortsByDistance(t *testing.T) {
	repo, db := setupUserRepository(t)
	ctx := context.Background()

	centerLat := 55.751244
	centerLon := 37.618423
	now := time.Now()

	farOwnerID := uuid.New()
	farLat := centerLat + 0.002
	farLon := centerLon
	require.NoError(t, db.Create(&user.User{
		ID:         farOwnerID,
		Email:      "sort-far@gav.app",
		Password:   "hash",
		Role:       "user",
		Visibility: user.VisibilityEveryone,
	}).Error)
	require.NoError(t, db.Create(&user.WalkSession{
		ID:         uuid.New(),
		UserID:     farOwnerID,
		Lat:        farLat,
		Lon:        farLon,
		Visibility: user.VisibilityEveryone,
		StartedAt:  now,
		UpdatedAt:  now,
	}).Error)
	require.NoError(t, db.Create(&dog.Dog{
		ID:       uuid.New(),
		OwnerID:  farOwnerID,
		Name:     "Дальше",
		Breed:    "Хаски",
		PhotoUrl: "/uploads/dogs/far.jpg",
		Status:   dog.StatusFriendly,
		Age:      dog.AdultAge,
		Gender:   dog.Male,
	}).Error)

	nearOwnerID := uuid.New()
	nearLat := centerLat + 0.001
	nearLon := centerLon
	require.NoError(t, db.Create(&user.User{
		ID:         nearOwnerID,
		Email:      "sort-near@gav.app",
		Password:   "hash",
		Role:       "user",
		Visibility: user.VisibilityEveryone,
	}).Error)
	require.NoError(t, db.Create(&user.WalkSession{
		ID:         uuid.New(),
		UserID:     nearOwnerID,
		Lat:        nearLat,
		Lon:        nearLon,
		Visibility: user.VisibilityEveryone,
		StartedAt:  now,
		UpdatedAt:  now,
	}).Error)
	require.NoError(t, db.Create(&dog.Dog{
		ID:       uuid.New(),
		OwnerID:  nearOwnerID,
		Name:     "Ближе",
		Breed:    "Корги",
		PhotoUrl: "/uploads/dogs/near.jpg",
		Status:   dog.StatusFriendly,
		Age:      dog.AdultAge,
		Gender:   dog.Female,
	}).Error)

	dogs, err := repo.FindWalkingNearby(ctx, centerLat, centerLon, 1_000, time.Now().Add(-2*time.Hour))

	require.NoError(t, err)
	require.Len(t, dogs, 2)
	require.Equal(t, "Ближе", dogs[0].Name)
	require.Equal(t, "Дальше", dogs[1].Name)
	require.Less(t, *dogs[0].DistanceMeters, *dogs[1].DistanceMeters)
}

func TestUserRepository_UpsertAndEndActiveWalkSession(t *testing.T) {
	repo, db := setupUserRepository(t)
	ctx := context.Background()
	userID := uuid.New()
	now := time.Now()

	require.NoError(t, repo.UpsertActiveWalkSession(ctx, &user.WalkSession{
		ID:         uuid.New(),
		UserID:     userID,
		Lat:        55.75,
		Lon:        37.61,
		Visibility: user.VisibilityEveryone,
		StartedAt:  now,
		UpdatedAt:  now,
	}))
	require.NoError(t, repo.UpsertActiveWalkSession(ctx, &user.WalkSession{
		ID:         uuid.New(),
		UserID:     userID,
		Lat:        55.76,
		Lon:        37.62,
		Visibility: user.VisibilityEveryone,
		StartedAt:  now.Add(time.Minute),
		UpdatedAt:  now.Add(time.Minute),
	}))

	var sessions []user.WalkSession
	require.NoError(t, db.Find(&sessions).Error)
	require.Len(t, sessions, 1)
	require.Equal(t, 55.76, sessions[0].Lat)
	require.Nil(t, sessions[0].EndedAt)

	require.NoError(t, repo.EndActiveWalkSession(ctx, userID, now.Add(2*time.Minute)))
	require.NoError(t, db.Find(&sessions).Error)
	require.Len(t, sessions, 1)
	require.NotNil(t, sessions[0].EndedAt)
}
