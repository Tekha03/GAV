package postgres

import (
	"context"
	"testing"

	"social_network/internal/stats"
	"social_network/internal/testdb"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func setupStatsRepository(t *testing.T) *StatsRepository {
	t.Helper()

	db := testdb.Open(t)
	require.NoError(t, db.AutoMigrate(&stats.UserStats{}, &stats.PostStats{}))

	repo, err := NewStatsRepository(db)
	require.NoError(t, err)

	return repo.(*StatsRepository)
}

func TestStatsRepository_IncrementUserFieldsCreatesMissingStats(t *testing.T) {
	repo := setupStatsRepository(t)
	ctx := context.Background()
	userID := uuid.New()

	require.NoError(t, repo.IncrementDogs(ctx, userID))
	require.NoError(t, repo.IncrementPosts(ctx, userID))

	got, err := repo.GetUserStats(ctx, userID)
	require.NoError(t, err)
	require.EqualValues(t, 1, got.DogsCount)
	require.EqualValues(t, 1, got.PostCount)
}

func TestStatsRepository_IncrementPostFieldsCreatesMissingStats(t *testing.T) {
	repo := setupStatsRepository(t)
	ctx := context.Background()
	postID := uuid.New()

	require.NoError(t, repo.IncrementPostLikes(ctx, postID))
	require.NoError(t, repo.IncrementPostComments(ctx, postID))

	got, err := repo.GetPostStats(ctx, postID)
	require.NoError(t, err)
	require.EqualValues(t, 1, got.LikesCount)
	require.EqualValues(t, 1, got.CommentsCount)
}
