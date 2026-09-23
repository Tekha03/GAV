package notification

import apperrors "shared/app_errors"

var (
	ErrEmptyHub              = apperrors.New(apperrors.Internal, "notification service: empty hub")
	ErrFailedToMarshal       = apperrors.New(apperrors.Internal, "failed to marshal notification")
	ErrNotificationRepoEmpty = apperrors.New(apperrors.NotificationNotFound, "notifications not found")
	ErrFirebaseClientEmpty   = apperrors.New(apperrors.Internal, "firebase client is empty")
)
