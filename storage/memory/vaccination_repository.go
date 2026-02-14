package memory

import (
	"context"
	"errors"
	"gav/internal/vaccination"
	"sync"

	"github.com/google/uuid"
)

var (
	ErrVaccinationExists = errors.New("vaccination already exists")
	ErrVaccinationNotFound = errors.New("vaccination not found")
)

type VaccinationRepository struct {
	mu 		sync.RWMutex
	vacs 	map[uuid.UUID]*vaccination.Vaccination
}

func NewVaccinationRepository() *VaccinationRepository {
	return &VaccinationRepository{
		vacs: make(map[uuid.UUID]*vaccination.Vaccination),
	}
}

func (r *VaccinationRepository) Create(ctx context.Context, v *vaccination.Vaccination) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if v.ID != uuid.Nil {
		if _, found := r.vacs[v.ID]; found {
			return ErrVaccinationExists
		}
	} else {
		v.ID = uuid.New()
	}

	if _, found := r.vacs[v.ID]; found {
		return ErrVaccinationExists
	}

	r.vacs[v.ID] = v
	return nil
}

func (r *VaccinationRepository) Update(ctx context.Context, v *vaccination.Vaccination) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, found := r.vacs[v.ID]; !found {
		return ErrVaccinationNotFound
	}

	r.vacs[v.ID] = v
	return nil
}

func (r *VaccinationRepository) Delete(ctx context.Context, ID uuid.UUID) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, found := r.vacs[ID]; found {
		return ErrVaccinationNotFound
	}

	delete(r.vacs, ID)
	return nil
}

func (r *VaccinationRepository) GetByID(ctx context.Context, ID uuid.UUID) (*vaccination.Vaccination, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, found := r.vacs[ID]; !found {
		return nil, ErrVaccinationNotFound
	}

	return r.vacs[ID], nil
}

func (r *VaccinationRepository) GetByDogID(ctx context.Context, dogID uuid.UUID) ([]*vaccination.Vaccination, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	var resultVacs []*vaccination.Vaccination
	for _, vac := range r.vacs {
		if vac.DogID == dogID {
			resultVacs = append(resultVacs, vac)
		}
	}

	return resultVacs, nil
}
