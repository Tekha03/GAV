package dog

import "context"


type DogRepository interface {
	Create(ctx context.Context, dog *Dog) error
	Update(ctx context.Context, dog *Dog) error
	Delete(ctx context.Context, ID uint) error
	GetByOwnerID(ctx context.Context, ownerID uint) ([]*Dog, error)
	GetByID(ctx context.Context, ID uint) (*Dog, error)
}