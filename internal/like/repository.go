package like

import "context"

type LikeRepository interface {
	Add(ctx context.Context, userID, postID uint) error
}
