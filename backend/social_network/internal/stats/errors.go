package stats

import apperrors "shared/app_errors"

var (
	ErrStatsNotFound = apperrors.New(apperrors.StatsNotFound, "user stats not found")
	ErrRepoNil       = apperrors.New(apperrors.Internal, "stats service: repo is nil")
)
