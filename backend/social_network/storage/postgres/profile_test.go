package postgres

import (
	"context"
	"testing"

	"social_network/internal/profile"
	"social_network/internal/testdb"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func setupProfileRepository(t *testing.T) *ProfileRepository {
	t.Helper()

	db := testdb.Open(t)
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
