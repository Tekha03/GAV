package memory

import (
	"context"
	"errors"
	"sync"

	"gav/internal/dog"
)

var (
	ErrDognotFound = errors.New("dog not found")
	ErrDogExists = errors.New("dog exists in repository")
)

type DogRepository struct {
	mu sync.RWMutex
	lastId	uint
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

	r.dogs[d.ID] = d
	return nil
}
