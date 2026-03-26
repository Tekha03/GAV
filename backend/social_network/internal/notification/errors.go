package notification

import "errors"

var (
	ErrEmptyHub			= errors.New("notification service: empty hub")
	ErrFailedToMarshal	= errors.New("failed to marshal into json")
	ErrNotificationRepoEmpty = errors.New("notifications not found")
	ErrFirebaseClientEmpty = errors.New("firebase client is empty")
)