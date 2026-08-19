package vaccination

import apperrors "shared/app_errors"

var (
	ErrVaccAccessDenied    = apperrors.New(apperrors.VaccinationAccessDenied, "vaccination access denied")
	ErrDogIDEmpty          = apperrors.New(apperrors.Validation, "dog ID cannot be empty", apperrors.WithDetail("field", "dog_id"))
	ErrDBError             = apperrors.New(apperrors.Internal, "database error")
	ErrRepoNil             = apperrors.New(apperrors.Internal, "vaccination service: repo is nil")
	ErrVaccinationNotFound = apperrors.New(apperrors.VaccinationNotFound, "vaccination not found")
)
