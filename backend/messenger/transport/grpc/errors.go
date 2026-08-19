package grpc

import (
	apperrors "shared/app_errors"

	"github.com/google/uuid"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func toGRPCError(err error) error {
	if err == nil {
		return nil
	}

	appErr := apperrors.Normalize(err)
	message := appErr.Message
	if appErr.Category == apperrors.CategoryInternal {
		message = "internal server error"
	}
	grpcStatus := status.New(grpcCode(appErr.Category), message)
	withDetails, detailsErr := grpcStatus.WithDetails(&errdetails.ErrorInfo{
		Reason: string(appErr.Code),
		Domain: "messenger",
	})
	if detailsErr != nil {
		return grpcStatus.Err()
	}
	return withDetails.Err()
}

func grpcCode(category apperrors.Category) codes.Code {
	switch category {
	case apperrors.CategoryValidation:
		return codes.InvalidArgument
	case apperrors.CategoryUnauthenticated:
		return codes.Unauthenticated
	case apperrors.CategoryPermissionDenied:
		return codes.PermissionDenied
	case apperrors.CategoryNotFound:
		return codes.NotFound
	case apperrors.CategoryConflict:
		return codes.AlreadyExists
	case apperrors.CategoryUnavailable:
		return codes.Unavailable
	case apperrors.CategoryUnsupported:
		return codes.Unimplemented
	case apperrors.CategoryCancelled:
		return codes.Canceled
	default:
		return codes.Internal
	}
}

func parseUUID(value, field string) (uuid.UUID, error) {
	id, err := uuid.Parse(value)
	if err != nil {
		return uuid.Nil, apperrors.Wrap(
			apperrors.Validation,
			"invalid UUID",
			err,
			apperrors.WithDetail("field", field),
		)
	}
	return id, nil
}

func requireCurrentUserID(ctxUserID uuid.UUID, ok bool) (uuid.UUID, error) {
	if !ok || ctxUserID == uuid.Nil {
		return uuid.Nil, apperrors.New(apperrors.AuthTokenMissing, "authenticated user is missing")
	}
	return ctxUserID, nil
}
