package postgres

import (
	"context"
	"errors"
	"sync"
	"testing"

	"social_network/internal/dog"
	"social_network/internal/follow"
	"social_network/internal/testdb"
	"social_network/internal/user"
	"social_network/internal/walk"

	"github.com/stretchr/testify/require"
)

// Run with GAV_TEST_POSTGRES_DSN against a disposable database.
func TestWalkPostgresPrivacyAndConcurrentStart(t *testing.T) {
	db := testdb.Open(t)
	require.NoError(t, db.AutoMigrate(&user.User{}, &user.WalkSession{}, &follow.Follow{}, &dog.Dog{}))
	require.NoError(t, db.Exec("CREATE UNIQUE INDEX IF NOT EXISTS idx_walk_sessions_active_user ON walk_sessions(user_id) WHERE ended_at IS NULL").Error)
	service := walk.NewService(NewWalkRepository(db))
	owner := createWalkUser(t, db)
	follower := createWalkUser(t, db)
	outsider := createWalkUser(t, db)

	var group sync.WaitGroup
	results := make(chan error, 2)
	for range 2 {
		group.Add(1)
		go func() {
			defer group.Done()
			_, err := service.Start(context.Background(), owner, 55.75, 37.61, user.VisibilityFollowersOnly)
			results <- err
		}()
	}
	group.Wait()
	close(results)
	successes, conflicts := 0, 0
	for err := range results {
		switch {
		case err == nil:
			successes++
		case errors.Is(err, walk.ErrAlreadyActive):
			conflicts++
		default:
			t.Fatalf("unexpected start error: %v", err)
		}
	}
	require.Equal(t, 1, successes)
	require.Equal(t, 1, conflicts)

	items, err := service.Nearby(context.Background(), outsider, 55.75, 37.61, 1000)
	require.NoError(t, err)
	require.Empty(t, items)
	require.NoError(t, db.Create(&follow.Follow{FollowerID: follower, FollowingID: owner}).Error)
	items, err = service.Nearby(context.Background(), follower, 55.75, 37.61, 1000)
	require.NoError(t, err)
	require.Len(t, items, 1)
	require.NoError(t, db.Where("follower_id = ? AND following_id = ?", follower, owner).Delete(&follow.Follow{}).Error)
	items, err = service.Nearby(context.Background(), follower, 55.75, 37.61, 1000)
	require.NoError(t, err)
	require.Empty(t, items)
}
