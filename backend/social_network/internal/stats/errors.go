package stats

import "errors"

var (
	ErrStatsNotFound = errors.New("user stats not found")
	ErrRepoNil       = errors.New("stats service: repo is nil")
)
