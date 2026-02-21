package memory

import (
	"context"
	"errors"
	"sync"

	"gav/internal/post"

	"github.com/google/uuid"
)

var (
	ErrPostNotFound = errors.New("post not found")
	ErrPostExists = errors.New("post already exists")
)

type PostRepository struct {
	mu sync.RWMutex
	posts map[uuid.UUID]*post.Post
}

func NewPostRepository() *PostRepository {
	return &PostRepository{
		posts: make(map[uuid.UUID]*post.Post),
	}
}

func (pr *PostRepository) Create(ctx context.Context, post *post.Post) error {
	pr.mu.Lock()
	defer pr.mu.Unlock()

	if _, found := pr.posts[post.ID]; found {
		return ErrPostExists
	}

	pr.posts[post.ID] = post
	return nil
}

func (pr *PostRepository) GetByID(ctx context.Context, id uuid.UUID) (*post.Post, error) {
	pr.mu.Lock()
	defer pr.mu.Unlock()

	foundPost, isOk := pr.posts[id]
	if !isOk {
		return nil, ErrPostNotFound
	}

	return foundPost, nil
}

func (pr *PostRepository) ListByUser(ctx context.Context, userID uuid.UUID) ([]*post.Post, error) {
	pr.mu.Lock()
	defer pr.mu.Unlock()

	var result []*post.Post
	for _, post := range pr.posts {
		if post.UserID == userID {
			result = append(result, post)
		}
	}

	return result, nil
}

func (pr *PostRepository) Delete(ctx context.Context, id uuid.UUID) error {
	pr.mu.Lock()
	defer pr.mu.Unlock()

	if _, isOk := pr.posts[id]; !isOk {
		return ErrPostNotFound
	}

	delete(pr.posts, id)
	return nil
}
