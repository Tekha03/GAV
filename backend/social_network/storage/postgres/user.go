package postgres

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
	err := r.DB(ctx).Where("settings_email = ?", u.Email).First(&existing).Error

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
	if err := r.DB(ctx).Where("settings_email = ?", email).First(&u).Error; err != nil {
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
	SELECT dogs.*
	FROM dogs
	JOIN users ON dogs.owner_id = users.id
	WHERE users.lat IS NOT NULL
		AND users.lon IS NOT NULL
		AND users.visibility = ?
		AND (
	  		6371000 * acos(
				LEAST(1.0, GREATEST(-1.0,
					cos(radians(?)) * cos(radians(users.lat)) * cos(radians(users.lon) - radians(?)) +
					sin(radians(?)) * sin(radians(users.lat))
				))
			)
		) <= ?
	`

	if err := r.DB(ctx).Raw(query, centerLat, centerLon, centerLat, radiusMeters).Scan(&dogs).Error; err != nil {
		return nil, fmt.Errorf("dog repository: find walking nearby: %w", err)
	}

	return dogs, nil
}
