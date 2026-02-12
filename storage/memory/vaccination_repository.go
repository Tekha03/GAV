package memory

import (
	"context"
	"errors"
	"gav/internal/vaccination"
	"sync"
)

var (
	ErrVaccinationExists   = errors.New("vaccination already exists")
	ErrVaccinationNotFound = errors.New("vaccination not found")
)

type VaccinationRepository struct {
	mu     sync.RWMutex
	lastID uint
	vacs   map[uint]*vaccination.Vaccination
}

func NewVaccinationRepository() *VaccinationRepository {
	return &VaccinationRepository{
		vacs: make(map[uint]*vaccination.Vaccination),
	}
}

func (r *VaccinationRepository) Create(ctx context.Context, v *vaccination.Vaccination) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.lastID++
	v.ID = r.lastID

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

func (r *VaccinationRepository) Delete(ctx context.Context, ID uint) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, found := r.vacs[ID]; found {
		return ErrVaccinationNotFound
	}

	delete(r.vacs, ID)
	return nil
}

func (r *VaccinationRepository) GetByID(ctx context.Context, ID uint) (*vaccination.Vaccination, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	vac, found := r.vacs[ID]
	if !found {
		return nil, ErrVaccinationNotFound
	}

	return vac, nil
}

func (r *VaccinationRepository) GetByDogID(ctx context.Context, dogID uint) ([]*vaccination.Vaccination, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var resultVacs []*vaccination.Vaccination
	for _, vac := range r.vacs {
		if vac.DogID == dogID {
			resultVacs = append(resultVacs, vac)
		}
	}

	return resultVacs, nil
}
