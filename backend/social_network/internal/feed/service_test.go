package feed

import (
	"context"
	"social_network/internal/post"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type testEnv struct {
	service FeedService
	repo    *MockPostRepository
}

func setup(t *testing.T) *testEnv {
	repo := &MockPostRepository{}

	service, err := NewService(repo)

	require.NoError(t, err)

	return &testEnv{
		service: service,
		repo:    repo,
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
}

func TestGetFeed_Success_NoNextPage(t *testing.T) {
	env := setup(t)

	ctx := context.Background()
	userID := uuid.New()
	before := time.Now()

	posts := []*post.Post{
		{CreatedAt: time.Now()},
		{CreatedAt: time.Now()},
	}

	env.repo.
		On("ListFeed", ctx, userID, before, 21).
		Return(posts, nil)

	result, next, err := env.service.GetFeed(ctx, userID, before, 10)

	assert.NoError(t, err)
	assert.Len(t, result, 2)
	assert.True(t, next.IsZero())
}

func TestGetFeed_LimitLessThanMin(t *testing.T) {
	env := setup(t)

	ctx := context.Background()
	userID := uuid.New()
	before := time.Now()

	env.repo.
		On("ListFeed", ctx, userID, before, 21).
		Return([]*post.Post{}, nil)

	_, _, err := env.service.GetFeed(ctx, userID, before, 5)

	assert.NoError(t, err)
}

func TestGetFeed_LimitGreaterThanMax(t *testing.T) {
	env := setup(t)

	ctx := context.Background()
	userID := uuid.New()
	before := time.Now()

	env.repo.
		On("ListFeed", ctx, userID, before, 101).
		Return([]*post.Post{}, nil)

	_, _, err := env.service.GetFeed(ctx, userID, before, 200)

	assert.NoError(t, err)
}

func TestGetFeed_WithNextPage(t *testing.T) {
	env := setup(t)

	ctx := context.Background()
	userID := uuid.New()
	before := time.Now()

	baseTime := time.Now()

	posts := make([]*post.Post, 21)
	for i := range posts {
		posts[i] = &post.Post{
			CreatedAt: baseTime.Add(time.Duration(i) * time.Minute),
		}
	}

	env.repo.
		On("ListFeed", ctx, userID, before, 21).
		Return(posts, nil)

	result, next, err := env.service.GetFeed(ctx, userID, before, 20)

	assert.NoError(t, err)
	assert.Len(t, result, 20)
	assert.Equal(t, posts[20].CreatedAt, next)
}

func TestGetFeed_RepoError(t *testing.T) {
	env := setup(t)

	ctx := context.Background()
	userID := uuid.New()
	before := time.Now()

	env.repo.
		On("ListFeed", ctx, userID, before, 21).
		Return(nil, assert.AnError)

	result, next, err := env.service.GetFeed(ctx, userID, before, 10)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.True(t, next.IsZero())
}

func TestGetFeed_ExactLimitPlusOne(t *testing.T) {
	env := setup(t)

	ctx := context.Background()
	userID := uuid.New()
	before := time.Now()

	baseTime := time.Now()

	posts := make([]*post.Post, 21)
	for i := range posts {
		posts[i] = &post.Post{
			CreatedAt: baseTime.Add(time.Duration(i)),
		}
	}

	env.repo.
		On("ListFeed", ctx, userID, before, 21).
		Return(posts, nil)

	result, next, err := env.service.GetFeed(ctx, userID, before, 20)

	assert.NoError(t, err)
	assert.Len(t, result, 20)
	assert.Equal(t, posts[20].CreatedAt, next)
}
