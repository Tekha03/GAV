package dog

type DogRepository interface {
	GetByOwnerID(ownerID uint) ([]*Dog, error)
	GetByID(id uint) (*Dog, error)
	Create(dog *Dog) error
	Update(dog *Dog) error
	Delete(id uint) error
}
