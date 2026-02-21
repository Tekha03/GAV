package sqlite

import (
	"context"
	"errors"
	"fmt"

	"gav/internal/dog"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type DogRepository struct {
	*BaseRepository
}

func NewDogRepository(db *gorm.DB) dog.Repository {
	return &DogRepository{BaseRepository: NewBaseRepository(db)}
}

func (r *DogRepository) Create(ctx context.Context, d *dog.Dog) error {
	if err := r.DB(ctx).Create(d).Error; err != nil {
		return fmt.Errorf("dog repository: create: %w", err)
	}

	return nil
}

func (r *DogRepository) Update(ctx context.Context, d *dog.Dog) error {
	result := r.DB(ctx).Model(&dog.Dog{}).Where("id = ?", d.ID).Updates(d)

	if result.Error != nil {
		return fmt.Errorf("dog repository: update: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return ErrDogNotFound
	}

	return nil
}

func (r *DogRepository) Delete(ctx context.Context, id uuid.UUID) error {
	result := r.DB(ctx).Delete(&dog.Dog{}, "id = ?", id)

	if result.Error != nil {
		return fmt.Errorf("dog repository: delete: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return ErrDogNotFound
	}

	return nil
}

func (r *DogRepository) GetByOwnerID(ctx context.Context, ownerID uuid.UUID) ([]*dog.Dog, error) {
	var dogs []*dog.Dog

	if err := r.DB(ctx).Where("owner_id = ?", ownerID).Find(&dogs).Error; err != nil {
		return nil, fmt.Errorf("dog repository: get by owner id: %w", err)
	}

	return dogs, nil
}

func (r *DogRepository) GetByID(ctx context.Context, id uuid.UUID) (*dog.Dog, error) {
	var d dog.Dog

	if err := r.DB(ctx).First(&d, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrDogNotFound
		}

		return nil, fmt.Errorf("dog repository: get by id: %w", err)
	}

	return &d, nil
}
