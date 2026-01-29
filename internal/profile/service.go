package profile

type ProfileService interface {
    Create(userID uint, input CreateProfileInput) (*UserProfile, error)
    Get(userID uint) (*UserProfile, error)
    Update(userID uint, input UpdateProfileInput) error
    Delete(userID uint) error
}
