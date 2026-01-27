package vaccination

type VaccinationRepository interface {
	Create(v Vaccination) (*Vaccination, error)
	Update(v Vaccination) error
	Delete(v Vaccination) error
	GetByDogID(dogid uint) ([]Vaccination, error)
	GetByID(vaccinationId uint) (*Vaccination, error)
}