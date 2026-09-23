package redis

import (
	"context"
	"errors"
	apperrors "shared/app_errors"
)

func unavailable(operation string, err error) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
		return apperrors.Normalize(err)
	}
	return apperrors.Wrap(apperrors.ServiceUnavailable, operation, err)
}

func internal(operation string, err error) error {
	if err == nil {
		return nil
	}
	return apperrors.Wrap(apperrors.Internal, operation, err)
}
