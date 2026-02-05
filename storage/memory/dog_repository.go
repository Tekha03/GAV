package memory

import (
	"context"
	"errors"
	"sync"

	"gav/internal/dog"
)

var (
	ErrDogNotFound = errors.New("dog not found")
	ErrDogExists = errors.New("dog exists in repository")
)

type DogRepository struct {
	mu 		sync.RWMutex
	lastID	uint
	dogs 	map[uint]*dog.Dog
}

func NewDogRepository() *DogRepository {
	return &DogRepository{
		dogs: make(map[uint]*dog.Dog),
	}
}

func (r *DogRepository) Create(ctx context.Context, d *dog.Dog) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, found := r.dogs[d.ID]; found {
		return ErrDogExists
	}

	d.ID = r.lastID
	r.dogs[d.ID] = d
	r.lastID++
	return nil
}

func (r *DogRepository) Update(ctx context.Context, d *dog.Dog) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, found := r.dogs[d.ID]; !found {
		return ErrDogNotFound
	}

	r.dogs[d.ID] = d
	return nil
}

func (r *DogRepository) GetByID(ctx context.Context, ID uint) (*dog.Dog, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, found := r.dogs[ID]; !found {
		return nil, ErrDogNotFound
	}

	d := r.dogs[ID]
	return d, nil
}

func (r *DogRepository) GetByOwnerID(ctx context.Context, ownerID uint) ([]*dog.Dog, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	var dogs []*dog.Dog
	for _, dog := range r.dogs {
		if dog.OwnerID == ownerID {
			dogs = append(dogs, dog)
		}
	}

	return dogs, nil
}

func (r *DogRepository) Delete(ctx context.Context, ID uint) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, found := r.dogs[ID]; !found {
		return ErrDogNotFound
	}

	delete(r.dogs, ID)
	return nil
}