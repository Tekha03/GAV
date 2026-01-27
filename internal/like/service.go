package like

import (
	"context"
)

type LikeService struct {
	repo LikeRepository
}

func (ls *LikeService) Add(ctx context.Context, userID, postID uint) error {
	return ls.repo.Add(ctx, userID, postID)
}
