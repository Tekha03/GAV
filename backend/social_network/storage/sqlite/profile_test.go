package sqlite

import (
	"context"
	"testing"

	"social_network/internal/profile"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupProfileRepository(t *testing.T) *ProfileRepository {
	t.Helper()

	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	require.NoError(t, err)
	require.NoError(t, db.AutoMigrate(&profile.UserProfile{}))

	repo, err := NewProfileRepository(db)
	require.NoError(t, err)

	return repo.(*ProfileRepository)
}

func TestProfileRepository_SearchFindsUsernameWithAtPrefix(t *testing.T) {
	repo := setupProfileRepository(t)
	ctx := context.Background()

	userProfile := &profile.UserProfile{
		UserID:   uuid.New(),
		Name:     "Виктория",
		Surname:  "К",
		Username: "vika_gav",
		Bio:      "Гуляем с собакой",
	}
	require.NoError(t, repo.Create(ctx, userProfile))

	profiles, err := repo.Search(ctx, "@vika", 10)

	require.NoError(t, err)
	require.Len(t, profiles, 1)
	require.Equal(t, userProfile.UserID, profiles[0].UserID)
}
