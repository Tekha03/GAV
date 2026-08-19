package dog

import apperrors "shared/app_errors"

var (
	ErrOwnerIDNil      = apperrors.New(apperrors.Validation, "dog owner id is required", apperrors.WithDetail("field", "owner_id"))
	ErrNameEmpty       = apperrors.New(apperrors.DogNameRequired, "dog name is required")
	ErrBreedEmpty      = apperrors.New(apperrors.DogBreedRequired, "dog breed is required")
	ErrGenderEmpty     = apperrors.New(apperrors.DogGenderInvalid, "dog gender is required")
	ErrStatusEmpty     = apperrors.New(apperrors.DogStatusInvalid, "dog status is required")
	ErrAgeEmpty        = apperrors.New(apperrors.DogAgeInvalid, "dog age is required")
	ErrPhotoURLEmpty   = apperrors.New(apperrors.Validation, "dog photo URL is required", apperrors.WithDetail("field", "photo_url"))
	ErrRepoNil         = apperrors.New(apperrors.Internal, "dog service: repo is nil")
	ErrDogAccessDenied = apperrors.New(apperrors.DogAccessDenied, "dog access denied")
)
