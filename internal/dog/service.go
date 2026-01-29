package dog

type DogService interface {
	Create(ownerID uint, input CreateDogInput) (*Dog, error)
    Update(ownerID, dogID uint, input UpdateDogInput) error
    Delete(ownerID, dogID uint) error

    UpdateLocation(ownerID, dogID uint, lat, lon float64) error
    SetLocationVisibility(ownerID, dogID uint, visible bool) error

    GetPublic(dogID uint) (*Dog, error)
    GetPrivate(ownerID, dogID uint) (*Dog, error)

	// later for analytics
	// GetStatusHistory(ownerID uint, dogID uint) ([]StatusChange, error)
}