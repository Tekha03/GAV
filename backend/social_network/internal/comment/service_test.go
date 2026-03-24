package comment

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type testEnv struct {
	service CommentService
	repo	*MockRepository
}

func setup(t *testing.T) *testEnv {
	repo := &MockRepository{}

	service, err := NewService(repo)
	require.NoError(t, err)

	return &testEnv{
		service: service,
		repo: 	 repo,
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
	assert.Equal(t, ErrRepoEmpty, err)
}

func TestCreate_Success(t *testing.T) {
	env := setup(t)
	ctx := context.Background()

	userID := uuid.New()
	postID := uuid.New()

	env.repo.
		On("Create", ctx, mock.AnythingOfType("*comment.Comment")).
		Return(nil)

	err := env.service.Create(ctx, userID, postID, "hello")

	assert.NoError(t, err)

	env.repo.AssertCalled(t, "Create", ctx, mock.Anything)
}

func TestCreate_RepoError(t *testing.T) {
	env := setup(t)
	ctx := context.Background()

	userID := uuid.New()
	postID := uuid.New()

	env.repo.
		On("Create", ctx, mock.Anything).
		Return(ErrDB)

	err := env.service.Create(ctx, userID, postID, "hello")
	assert.Error(t, err)
}

func TestGetByID_Success(t *testing.T) {
	env := setup(t)
	ctx := context.Background()

	id := uuid.New()
	comment := &Comment{
		ID: id,
		Content: "text",
	}

	env.repo.
		On("GetByID", ctx, id).
		Return(comment, nil)

	result, err := env.service.GetByID(ctx, id)

	assert.NoError(t, err)
	assert.Equal(t, comment, result)
}

func TestGetByID_NotFound(t *testing.T) {
	env := setup(t)
	ctx := context.Background()

	id := uuid.New()
	env.repo.
		On("GetByID", ctx, id).
		Return(nil, ErrCommentNotFound)

	result, err := env.service.GetByID(ctx, id)

	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestListByPostID_Success(t *testing.T) {
	env := setup(t)
	ctx := context.Background()

	postID := uuid.New()
	comments := []Comment{
		{ID: uuid.New(), Content: "1"},
		{ID: uuid.New(), Content: "2"},
	}

	env.repo.
		On("ListByPostID", ctx, postID).
		Return(comments, nil)

	result, err := env.service.ListByPostID(ctx, postID)

	assert.NoError(t, err)
	assert.Len(t, result, 2)
}

func TestListByPostID_Error(t *testing.T) {
	env := setup(t)
	ctx := context.Background()

	postID := uuid.New()
	env.repo.
		On("ListByPostID", ctx, postID).
		Return(nil, ErrDB)

	result, err := env.service.ListByPostID(ctx, postID)

	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestDelete_Success(t *testing.T) {
	env := setup(t)
	ctx := context.Background()

	userID := uuid.New()
	commentID := uuid.New()

	env.repo.
		On("Delete", ctx, userID, commentID).
		Return(nil)

	err := env.service.Delete(ctx, userID, commentID)
	assert.NoError(t, err)
}

func TestDelete_Error(t *testing.T) {
	env := setup(t)
	ctx := context.Background()

	userID := uuid.New()
	commentID := uuid.New()

	env.repo.
		On("Delete", ctx, userID, commentID).
		Return(ErrDB)

	err := env.service.Delete(ctx, userID, commentID)
	assert.Error(t, err)
}