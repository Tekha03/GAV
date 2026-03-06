package sqlite

import (
	"context"
	"errors"
	"fmt"
	"gav/internal/vaccination"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type VaccinationRepository struct {
	*BaseRepository
}

func NewVaccinationRepository(db *gorm.DB) (vaccination.Repository, error) {
	repo, err := NewBaseRepository(db)
	if err != nil {
		return nil, err
	}

	return &VaccinationRepository{BaseRepository: repo}, nil
}

func (r *VaccinationRepository) Create(ctx context.Context, v *vaccination.Vaccination) error {
	if err := r.DB(ctx).Create(v).Error; err != nil {
		return fmt.Errorf("vaccination repository: create: %w", err)
	}

	return nil
}

func (r *VaccinationRepository) Update(ctx context.Context, v *vaccination.Vaccination) error {
	result := r.DB(ctx).Model(&vaccination.Vaccination{}).Where("id = ?", v.ID).Updates(v)

	if result.Error != nil {
		return fmt.Errorf("vaccination repository: update: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return ErrVaccinationNotFound
	}

	return nil
}

func (r *VaccinationRepository) Delete(ctx context.Context, id uuid.UUID) error {
	result := r.DB(ctx).Delete(&vaccination.Vaccination{}, "id = ?", id)

	if result.Error != nil {
		return fmt.Errorf("vaccination repository: delete: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return ErrVaccinationNotFound
	}

	return nil
}

func (r *VaccinationRepository) ListByDogID(ctx context.Context, dogID uuid.UUID) ([]*vaccination.Vaccination, error) {
	var vaccinations []*vaccination.Vaccination

	if err := r.DB(ctx).Where("dog_id = ?", dogID).Order("created_at DESC").Find(&vaccinations).Error; err != nil {
		return nil, fmt.Errorf("vaccination repository: list by dog id: %w", err)
	}

	return vaccinations, nil
}

func (r *VaccinationRepository) GetByID(ctx context.Context, id uuid.UUID) (*vaccination.Vaccination, error) {
	var v vaccination.Vaccination

	if err := r.DB(ctx).First(&v, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrVaccinationNotFound
		}

		return nil, fmt.Errorf("vaccination repository: get by id: %w", err)
	}

	return &v, nil
}
