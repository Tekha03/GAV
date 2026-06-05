package like

import "context"

type Repository interface {
	Add(ctx context.Context, like Like) error
	Remove(ctx context.Context, like Like) error
	LikeExists(ctx context.Context, like Like) (bool, error)
}
