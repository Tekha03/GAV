package postgres

import (
	"context"
	"errors"
	"math"
	"sync"
	"testing"
	"time"

	"social_network/internal/dog"
	"social_network/internal/follow"
	"social_network/internal/testdb"
	"social_network/internal/user"
	"social_network/internal/walk"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func setupWalk(t *testing.T) (*walk.Service, *gorm.DB) {
	t.Helper()
	db := testdb.Open(t)
	require.NoError(t, db.AutoMigrate(&user.User{}, &user.WalkSession{}, &follow.Follow{}, &dog.Dog{}))
	require.NoError(t, db.Exec("CREATE UNIQUE INDEX idx_walk_sessions_active_user ON walk_sessions(user_id) WHERE ended_at IS NULL").Error)
	return walk.NewService(NewWalkRepository(db)), db
}

func createWalkUser(t *testing.T, db *gorm.DB) uuid.UUID {
	t.Helper()
	id := uuid.New()
	require.NoError(t, db.Create(&user.User{ID: id, Email: id.String() + "@test.local", Password: "hash", Role: "user"}).Error)
	return id
}

func TestWalkLifecycleAndVisibility(t *testing.T) {
	service, db := setupWalk(t)
	ctx := context.Background()
	owner := createWalkUser(t, db)
	follower := createWalkUser(t, db)
	other := createWalkUser(t, db)
	const lat, lon = 55.751244, 37.618423

	_, err := service.Start(ctx, owner, lat, lon, user.VisibilityFollowersOnly)
	require.NoError(t, err)
	var account user.User
	require.NoError(t, db.First(&account, "id = ?", owner).Error)
	require.Equal(t, user.Walking, account.LocationStatus)
	require.Equal(t, user.VisibilityFollowersOnly, account.Visibility)
	_, err = service.Start(ctx, owner, lat, lon, user.VisibilityEveryone)
	require.ErrorIs(t, err, walk.ErrAlreadyActive)

	current, err := service.Current(ctx, owner)
	require.NoError(t, err)
	require.Equal(t, owner, current.UserID)
	results, err := service.Nearby(ctx, other, lat, lon, 1000)
	require.NoError(t, err)
	require.Empty(t, results)
	require.NoError(t, db.Create(&follow.Follow{FollowerID: owner, FollowingID: other}).Error)
	results, err = service.Nearby(ctx, other, lat, lon, 1000)
	require.NoError(t, err)
	require.Empty(t, results, "reverse follow must not grant access")

	// A follows B: A can see B; the reverse relationship is not enough.
	require.NoError(t, db.Create(&follow.Follow{FollowerID: follower, FollowingID: owner}).Error)
	results, err = service.Nearby(ctx, follower, lat, lon, 1000)
	require.NoError(t, err)
	require.Len(t, results, 1)
	require.Equal(t, owner, results[0].UserID)

	require.NoError(t, db.Where("follower_id = ? AND following_id = ?", follower, owner).Delete(&follow.Follow{}).Error)
	results, err = service.Nearby(ctx, follower, lat, lon, 1000)
	require.NoError(t, err)
	require.Empty(t, results)

	_, err = service.UpdateVisibility(ctx, owner, user.VisibilityEveryone)
	require.NoError(t, err)
	results, err = service.Nearby(ctx, other, lat, lon, 1000)
	require.NoError(t, err)
	require.Len(t, results, 1)

	_, err = service.UpdateVisibility(ctx, owner, user.VisibilityNoOne)
	require.NoError(t, err)
	results, err = service.Nearby(ctx, other, lat, lon, 1000)
	require.NoError(t, err)
	require.Empty(t, results)
	_, err = service.Current(ctx, owner)
	require.NoError(t, err)

	require.NoError(t, service.Stop(ctx, owner))
	require.NoError(t, db.First(&account, "id = ?", owner).Error)
	require.Equal(t, user.Inactive, account.LocationStatus)
	require.NoError(t, service.Stop(ctx, owner))
	_, err = service.Current(ctx, owner)
	require.ErrorIs(t, err, walk.ErrNotFound)
	_, err = service.UpdateLocation(ctx, owner, lat, lon)
	require.ErrorIs(t, err, walk.ErrNotFound)
}

func TestWalkStartRollsBackIfUserStateUpdateFails(t *testing.T) {
	service, db := setupWalk(t)
	owner := createWalkUser(t, db)
	require.NoError(t, db.Exec("CREATE FUNCTION reject_user_update() RETURNS trigger LANGUAGE plpgsql AS $$ BEGIN RAISE EXCEPTION 'reject'; END $$").Error)
	require.NoError(t, db.Exec("CREATE TRIGGER reject_user_update BEFORE UPDATE ON users FOR EACH ROW EXECUTE FUNCTION reject_user_update()").Error)
	_, err := service.Start(context.Background(), owner, 55.75, 37.61, user.VisibilityEveryone)
	require.Error(t, err)
	var count int64
	require.NoError(t, db.Model(&user.WalkSession{}).Where("user_id = ? AND ended_at IS NULL", owner).Count(&count).Error)
	require.Zero(t, count)
}

func TestWalkNearbyRadiusTTLAndOneResultPerOwner(t *testing.T) {
	service, db := setupWalk(t)
	ctx := context.Background()
	requester := createWalkUser(t, db)
	near := createWalkUser(t, db)
	far := createWalkUser(t, db)
	stale := createWalkUser(t, db)
	const lat, lon = 55.75, 37.61
	for _, item := range []struct {
		id       uuid.UUID
		latitude float64
	}{{near, lat + 0.001}, {far, lat + 0.02}, {stale, lat + 0.0005}} {
		_, err := service.Start(ctx, item.id, item.latitude, lon, user.VisibilityEveryone)
		require.NoError(t, err)
	}
	for range 2 {
		require.NoError(t, db.Create(&dog.Dog{ID: uuid.New(), OwnerID: near, Name: "dog", Breed: "breed"}).Error)
	}
	require.NoError(t, db.Model(&user.WalkSession{}).Where("user_id = ?", stale).Update("updated_at", time.Now().Add(-3*time.Hour)).Error)

	results, err := service.Nearby(ctx, requester, lat, lon, 500)
	require.NoError(t, err)
	require.Len(t, results, 1)
	require.Equal(t, near, results[0].UserID)
	require.Len(t, results[0].Dogs, 2)
	require.Equal(t, "dog", results[0].Dogs[0].Name)
	require.InDelta(t, 111, results[0].DistanceMeters, 5)
	_, err = service.Current(ctx, stale)
	require.True(t, errors.Is(err, walk.ErrNotFound))
}

func TestWalkStartRollsBackWhenInsertFails(t *testing.T) {
	service, db := setupWalk(t)
	ctx := context.Background()
	owner := createWalkUser(t, db)
	_, err := service.Start(ctx, owner, 55.75, 37.61, user.VisibilityEveryone)
	require.NoError(t, err)
	require.NoError(t, db.Model(&user.WalkSession{}).Where("user_id = ?", owner).Update("updated_at", time.Now().Add(-3*time.Hour)).Error)
	require.NoError(t, db.Exec("CREATE FUNCTION reject_walk_insert() RETURNS trigger LANGUAGE plpgsql AS $$ BEGIN RAISE EXCEPTION 'reject'; END $$").Error)
	require.NoError(t, db.Exec("CREATE TRIGGER reject_walk_insert BEFORE INSERT ON walk_sessions FOR EACH ROW EXECUTE FUNCTION reject_walk_insert()").Error)
	_, err = service.Start(ctx, owner, 55.75, 37.61, user.VisibilityEveryone)
	require.Error(t, err)
	var session user.WalkSession
	require.NoError(t, db.Where("user_id = ?", owner).First(&session).Error)
	require.Nil(t, session.EndedAt, "stale session must remain unchanged after rollback")
}

func TestWalkNearbyCrossesAntimeridian(t *testing.T) {
	service, db := setupWalk(t)
	ctx := context.Background()
	requester := createWalkUser(t, db)
	owner := createWalkUser(t, db)
	_, err := service.Start(ctx, owner, 0, -179.999, user.VisibilityEveryone)
	require.NoError(t, err)
	results, err := service.Nearby(ctx, requester, 0, 179.999, 500)
	require.NoError(t, err)
	require.Len(t, results, 1)
}

func TestWalkValidation(t *testing.T) {
	service, _ := setupWalk(t)
	_, err := service.Start(context.Background(), uuid.New(), 91, 0, user.VisibilityEveryone)
	require.ErrorIs(t, err, walk.ErrInvalidCoordinates)
	_, err = service.Start(context.Background(), uuid.New(), 0, 0, 99)
	require.ErrorIs(t, err, walk.ErrInvalidVisibility)
	_, err = service.Nearby(context.Background(), uuid.New(), 0, 0, 20)
	require.ErrorIs(t, err, walk.ErrInvalidRadius)
	_, err = service.Nearby(context.Background(), uuid.New(), math.NaN(), 0, 1000)
	require.ErrorIs(t, err, walk.ErrInvalidCoordinates)
	_, err = service.Nearby(context.Background(), uuid.New(), 0, 0, math.Inf(1))
	require.ErrorIs(t, err, walk.ErrInvalidRadius)
}

func TestWalkNearbySortsByDistance(t *testing.T) {
	service, db := setupWalk(t)
	ctx := context.Background()
	requester := createWalkUser(t, db)
	near := createWalkUser(t, db)
	far := createWalkUser(t, db)
	_, err := service.Start(ctx, far, 55.753, 37.61, user.VisibilityEveryone)
	require.NoError(t, err)
	_, err = service.Start(ctx, near, 55.751, 37.61, user.VisibilityEveryone)
	require.NoError(t, err)
	results, err := service.Nearby(ctx, requester, 55.75, 37.61, 1000)
	require.NoError(t, err)
	require.Len(t, results, 2)
	require.Equal(t, near, results[0].UserID)
	require.Equal(t, far, results[1].UserID)
	require.Less(t, results[0].DistanceMeters, results[1].DistanceMeters)
}

func TestConcurrentWalkStartCreatesOneActiveSession(t *testing.T) {
	service, db := setupWalk(t)
	connection, err := db.DB()
	require.NoError(t, err)
	connection.SetMaxOpenConns(1)
	owner := createWalkUser(t, db)
	var group sync.WaitGroup
	results := make(chan error, 2)
	for range 2 {
		group.Add(1)
		go func() {
			defer group.Done()
			_, err := service.Start(context.Background(), owner, 55.75, 37.61, user.VisibilityEveryone)
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
	var count int64
	require.NoError(t, db.Model(&user.WalkSession{}).Where("user_id = ? AND ended_at IS NULL", owner).Count(&count).Error)
	require.EqualValues(t, 1, count)
}
