package sqlite

import (
	"context"
	"errors"
	"fmt"

	"social_network/internal/dog"
	"social_network/internal/user"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserRepository struct {
	*BaseRepository
}

func NewUserRepository(db *gorm.DB) (user.Repository, error) {
	repo, err := NewBaseRepository(db)
	if err != nil {
		return nil, err
	}

	return &UserRepository{BaseRepository: repo}, nil
}

func (r *UserRepository) Create(ctx context.Context, u *user.User) error {
	var existing user.User
	err := r.DB(ctx).Where("email = ?", u.Email).First(&existing).Error

	if err == nil {
		return ErrUserExists
	}

	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	return r.DB(ctx).Create(u).Error
}

func (r *UserRepository) GetByID(ctx context.Context, id uuid.UUID) (*user.User, error) {
	var u user.User
	if err := r.DB(ctx).First(&u, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}

		return nil, err
	}

	return &u, nil
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*user.User, error) {
	var u user.User
	if err := r.DB(ctx).Where("email = ?", email).First(&u).Error; err != nil {
		return nil, err
	}

	return &u, nil
}

func (r *UserRepository) Update(ctx context.Context, u *user.User) error {
	updated := r.DB(ctx).Save(u)
	if updated.Error != nil {
		return updated.Error
	}

	if updated.RowsAffected == 0 {
		return ErrUserNotFound
	}

	return nil
}

func (r *UserRepository) Delete(ctx context.Context, id uuid.UUID) error {
	deleted := r.DB(ctx).Delete(&user.User{}, "id = ?", id)
	if deleted.Error != nil {
		return deleted.Error
	}

	if deleted.RowsAffected == 0 {
		return ErrUserNotFound
	}

	return nil
}

func (r *UserRepository) FindWalkingNearby(ctx context.Context, centerLat, centerLon float64, radiusMeters float64) ([]*dog.Dog, error) {
	var dogs []*dog.Dog

	query := `
	SELECT
		dogs.id,
		dogs.owner_id,
		dogs.name,
		dogs.breed,
		dogs.photo_url,
		dogs.status,
		dogs.age,
		dogs.gender,
		users.lat AS lat,
		users.lon AS lon
	FROM dogs
	JOIN users ON users.id = dogs.owner_id
	WHERE users.lat IS NOT NULL AND users.lon IS NOT NULL
	  AND users.visibility = ?
	  AND users.location_status = ?
	  AND (((users.lat - ?) * 111320) * ((users.lat - ?) * 111320)
	    + ((users.lon - ?) * 111320) * ((users.lon - ?) * 111320)) <= (? * ?)
	`

	if err := r.DB(ctx).Raw(
		query,
		user.VisibilityEveryone,
		user.Walking,
		centerLat,
		centerLat,
		centerLon,
		centerLon,
		radiusMeters,
		radiusMeters,
	).Scan(&dogs).Error; err != nil {
		return nil, fmt.Errorf("dog repository: find walking nearby: %w", err)
	}

	return dogs, nil
}
