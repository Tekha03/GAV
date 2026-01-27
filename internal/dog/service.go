package dog

import "gav/internal/vaccination"


type DogService interface {
	CreateDog(ownerID uint, dog *Dog) (*Dog, error)
	UpdateStatus(ownerID, dogID uint, newstatus Status) error
	UpdatePhoto(ownerID, dogID uint, photoUrl string) error
	UpdateAge(ownerID, dogID uint, newAge Age) error
	UpdateInfo(ownerID uint, dogID uint, name, breed, gender string) error
	DeleteDog(ownerID, dogID uint) error

	AddVaccination(ownerID, dogID uint, v vaccination.Vaccination) (*vaccination.Vaccination, error)
	UpdateVaccination(ownerID, dogID uint, v vaccination.Vaccination) error
	DeleteVaccination(ownerID, dogID, vaccinationID uint) error
	GetVaccinations(ownerID, dogID uint) ([]vaccination.Vaccination, error)

	UpdateLocation(ownerID, dogID uint, lat, lon float64) error
    SetLocationVisibility(ownerID, dogID uint, visible bool) error

	GetPublic(dogID uint) (*Dog, error)
    GetPrivate(ownerID, dogID uint) (*Dog, error)

	// later for analytics
	// GetStatusHistory(ownerID uint, dogID uint) ([]StatusChange, error)
}