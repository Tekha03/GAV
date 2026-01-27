package vaccination

type VaccinationService interface {
	AddVaccination(ownerID uint, dogID uint, v Vaccination) (*Vaccination, error)
	UpdateVaccination(ownerID uint, dogID uint, v Vaccination) error
	DeleteVaccination(ownerID uint, dogID uint, vaccinationID uint) error
	GetVaccinations(ownerID uint, dogID uint) ([]Vaccination, error)
}