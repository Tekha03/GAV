package dog

type DogService struct {
	repo 	DogRepository
}

func NewDogService(repo DogRepository) DogService {
	return DogService{repo: repo}
}

func (s *DogService) Create(ownerID uint, input CreateDogInput) (*Dog, error) {
	dog := NewDog(
		ownerID, 
		input.Name, 
		input.Breed, 
		input.Gender, 
		input.Status, 
		input.Age, 
		input.PhotoUrl,
	)

	if err := s.repo.Create(dog); err != nil {
		return nil, err
	}

	return dog, nil
}

func (d *Dog) Update(ownerID uint, input UpdateDogInput) error {
	if input.Name != nil {
		d.Name = *input.Name
	}

	if input.Breed != nil {
		d.Breed = *input.Breed
	}

	if input.PhotoUrl != nil {
		d.PhotoUrl = *input.PhotoUrl
	}

	if input.Age != nil {
		d.Age = *input.Age
	}

	if input.Gender != nil {
		d.Gender = *input.Gender
	}

	if input.Status != nil {
		d.Status = *input.Status
	}

	return nil
}