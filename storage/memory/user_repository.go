package memory

import (
	"context"
	"sync"

	"gav/internal/user"

	"github.com/google/uuid"
)

type UserRepository struct {
	mu		sync.RWMutex
	byID	map[uuid.UUID]*user.User
	byEmail	map[string]*user.User
}

func NewUserRepository() *UserRepository {
	return &UserRepository{
		byID: 	 make(map[uuid.UUID]*user.User),
		byEmail: make(map[string]*user.User),
	}
}

func (r *UserRepository) Create(ctx context.Context, u *user.User) error {
	if u == nil {
		return ErrUserNil
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.byEmail[u.Email]; exists {
		return ErrUserExists
	}

	r.byID[u.ID] = u
	r.byEmail[u.Email] = u
	return nil
}

func (r *UserRepository) GetByID(ctx context.Context, id uuid.UUID) (*user.User, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	foundUser, isOk := r.byID[id]
	if !isOk {
		return nil, ErrUserNotFound
	}

	return foundUser, nil
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*user.User, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	foundUser, isOk := r.byEmail[email]
	if !isOk {
		return nil, ErrUserNotFound
	}

	return foundUser, nil
}

func (r *UserRepository) Update(ctx context.Context, user *user.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, isOk := r.byID[user.ID]; !isOk {
		return ErrUserNotFound
	}

	r.byID[user.ID] = user
	r.byEmail[user.Email] = user
	return nil
}

func (r *UserRepository) Delete(ctx context.Context, id uuid.UUID) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	foundUser, isOk := r.byID[id]
	if !isOk {
		return ErrUserNotFound
	}

	delete(r.byID, id)
	delete(r.byEmail, foundUser.Email)
	return nil
}
