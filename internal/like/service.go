package like

import "context"

type LikeService interface {
	Add(ctx context.Context, like Like) error
	Remove(ctx context.Context, like Like) error
}
