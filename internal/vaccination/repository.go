package vaccination

type VaccinationRepository interface {
	Create(v *Vaccination) (*Vaccination, error)
	Update(v *Vaccination) error
	Delete(dogID uint) error
	GetByDogID(dogID uint) ([]Vaccination, error)
	GetByID(vaccinationId uint) (*Vaccination, error)
}