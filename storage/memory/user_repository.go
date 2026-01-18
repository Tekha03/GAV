package memory

import (
	"context"
	"errors"
	"sync"

	"gav/user"
)

var (
	ErrUserNotFound = errors.New("user not found")
	ErrUserExists = errors.New("user already exists")
)

type UserRepository struct {
	mu		sync.RWMutex
	byID	map[uint]*user.User
	byEmail	map[string]*user.User
}

func NewUserRepository() *UserRepository {
	return &UserRepository{
		byID: 	 make(map[uint]*user.User),
		byEmail: make(map[string]*user.User),
	}
}

func (ur *UserRepository) Create(ctx context.Context, u *user.User) error {
	ur.mu.Lock()
	defer ur.mu.Unlock()

	if _, exists := ur.byEmail[u.Profile.Email]; exists {
		return ErrUserExists
	}

	ur.byID[u.ID] = u
	ur.byEmail[u.Profile.Email] = u
	return nil
}

func (ur *UserRepository) GetByID(ctx context.Context, id uint) (*user.User, error) {
	ur.mu.Lock()
	defer ur.mu.Unlock()

	foundUser, isOk := ur.byID[id]
	if !isOk {
		return nil, ErrUserNotFound
	}

	return foundUser, nil
}

func (ur *UserRepository) GetByEmail(ctx context.Context, email string) (*user.User, error) {
	ur.mu.Lock()
	defer ur.mu.Unlock()

	foundUser, isOk := ur.byEmail[email]
	if !isOk {
		return nil, ErrUserNotFound
	}

	return foundUser, nil
}

func (ur *UserRepository) Update(ctx context.Context, user *user.User) error {
	ur.mu.Lock()
	defer ur.mu.Unlock()

	if _, isOk := ur.byID[user.ID]; !isOk {
		return ErrUserNotFound
	}

	ur.byID[user.ID] = user
	ur.byEmail[user.Profile.Email] = user
	return nil
}

func (ur *UserRepository) Delete(ctx context.Context, id uint) error {
	ur.mu.Lock()
	defer ur.mu.Unlock()

	foundUser, isOk := ur.byID[id]
	if !isOk {
		return ErrUserNotFound
	}

	delete(ur.byID, id)
	delete(ur.byEmail, foundUser.Profile.Email)
	return nil
}
