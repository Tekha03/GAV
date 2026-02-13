package dog

import "github.com/google/uuid"

type Service interface {
	Create(ownerID uuid.UUID, input CreateDogInput) (*Dog, error)
    Update(ownerID, dogID uuid.UUID, input UpdateDogInput) error
    Delete(ownerID, dogID uuid.UUID) error

    UpdateLocation(ownerID, dogID uuid.UUID, lat, lon float64) error
    SetLocationVisibility(ownerID, dogID uuid.UUID, visible bool) error

    GetPublic(dogID uuid.UUID) (*Dog, error)
    GetPrivate(ownerID, dogID, dogOwnerID uuid.UUID) (*Dog, error)

	// later for analytics
	// GetStatusHistory(ownerID uint, dogID uint) ([]StatusChange, error)
}
