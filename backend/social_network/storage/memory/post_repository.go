package memory

import (
	"context"
	"errors"
	"sync"

	"social_network/internal/post"

	"github.com/google/uuid"
)

var (
	ErrPostNotFound = errors.New("post not found")
	ErrPostExists   = errors.New("post already exists")
)

type PostRepository struct {
	mu    sync.RWMutex
	posts map[uuid.UUID]*post.Post
}

func NewPostRepository() *PostRepository {
	return &PostRepository{
		posts: make(map[uuid.UUID]*post.Post),
	}
}

func (r *PostRepository) Create(ctx context.Context, post *post.Post) error {
	if post == nil {
		return ErrPostNil
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	if _, found := r.posts[post.ID]; found {
		return ErrPostExists
	}

	r.posts[post.ID] = post
	return nil
}

func (r *PostRepository) GetByID(ctx context.Context, id uuid.UUID) (*post.Post, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	foundPost, isOk := r.posts[id]
	if !isOk {
		return nil, ErrPostNotFound
	}

	return foundPost, nil
}

func (r *PostRepository) ListByUser(ctx context.Context, userID uuid.UUID) ([]*post.Post, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	var result []*post.Post
	for _, post := range r.posts {
		if post.UserID == userID {
			result = append(result, post)
		}
	}

	return result, nil
}

func (r *PostRepository) Delete(ctx context.Context, id uuid.UUID) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, isOk := r.posts[id]; !isOk {
		return ErrPostNotFound
	}

	delete(r.posts, id)
	return nil
}
