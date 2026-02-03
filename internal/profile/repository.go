package profile

type ProfileRepository interface {
	Create(profile *UserProfile) error
	Update(profile *UserProfile) error
	Delete(profileID uint) error
	GetByID(profileID uint) (*UserProfile, error)
}
