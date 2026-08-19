package gorm

import (
	"context"
	"errors"
	apperrors "shared/app_errors"

	"gorm.io/gorm"
)

func createError(err error, conflict apperrors.Definition, conflictMessage, operation string) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, gorm.ErrDuplicatedKey) {
		return apperrors.Wrap(conflict, conflictMessage, err)
	}
	return internalError(operation, err)
}

func mutationError(result *gorm.DB, notFound apperrors.Definition, notFoundMessage, operation string) error {
	if result.Error != nil {
		return internalError(operation, result.Error)
	}
	if result.RowsAffected == 0 {
		return apperrors.New(notFound, notFoundMessage)
	}
	return nil
}

func internalError(operation string, err error) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
		return apperrors.Normalize(err)
	}
	return apperrors.Wrap(apperrors.Internal, operation, err)
}
