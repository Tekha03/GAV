package profile

import "context"

type ProfileService interface {
    Create(ctx context.Context, userID uint, input CreateProfileInput) (*UserProfile, error)
    GetByID(ctx context.Context, profileID uint) (*UserProfile, error)
    Update(ctx context.Context, profileID uint, input UpdateProfileInput) error
    Delete(ctx context.Context, profileID uint) error
}
