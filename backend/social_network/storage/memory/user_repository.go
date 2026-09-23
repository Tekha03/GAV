package memory

import (
	"context"
	"sync"
	"time"

	"social_network/internal/dog"
	"social_network/internal/user"

	"github.com/google/uuid"
)

type UserRepository struct {
	mu           sync.RWMutex
	byID         map[uuid.UUID]*user.User
	byEmail      map[string]*user.User
	walkSessions map[uuid.UUID]*user.WalkSession
}

func NewUserRepository() *UserRepository {
	return &UserRepository{
		byID:         make(map[uuid.UUID]*user.User),
		byEmail:      make(map[string]*user.User),
		walkSessions: make(map[uuid.UUID]*user.WalkSession),
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

func (r *UserRepository) FindWalkingNearby(ctx context.Context, centerLat, centerLon float64, radiusMeters float64, activeAfter time.Time) ([]*dog.Dog, error) {
	return []*dog.Dog{}, nil
}

func (r *UserRepository) UpsertActiveWalkSession(ctx context.Context, session *user.WalkSession) error {
	if session == nil {
		return nil
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	if existing := r.walkSessions[session.UserID]; existing != nil && existing.EndedAt == nil {
		existing.Lat = session.Lat
		existing.Lon = session.Lon
		existing.Visibility = session.Visibility
		existing.UpdatedAt = session.UpdatedAt
		return nil
	}

	r.walkSessions[session.UserID] = session
	return nil
}

func (r *UserRepository) EndActiveWalkSession(ctx context.Context, userID uuid.UUID, endedAt time.Time) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if session := r.walkSessions[userID]; session != nil && session.EndedAt == nil {
		session.EndedAt = &endedAt
		session.UpdatedAt = endedAt
	}

	return nil
}

func (r *UserRepository) UpdateActiveWalkVisibility(ctx context.Context, userID uuid.UUID, visibility user.LocationVisibility) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if session := r.walkSessions[userID]; session != nil && session.EndedAt == nil {
		session.Visibility = visibility
	}

	return nil
}
