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

func (s *DogService) Update(ownerID, dogID uint, input UpdateDogInput) error {

	dog, err := s.repo.GetByID(dogID)
    if err != nil {
        return err
    }

	if input.Name != nil {
		dog.Name = *input.Name
	}

	if input.Breed != nil {
		dog.Breed = *input.Breed
	}

	if input.PhotoUrl != nil {
		dog.PhotoUrl = *input.PhotoUrl
	}

	if input.Age != nil {
		dog.Age = *input.Age
	}

	if input.Gender != nil {
		dog.Gender = *input.Gender
	}

	if input.Status != nil {
		dog.Status = *input.Status
	}

	return s.repo.Update(dog)
}

func (s *DogService) GetPublic(dogID uint) (*Dog, error) {
	return s.repo.GetByID(dogID)
}

func (s *DogService) GetPrivate(dogID uint)