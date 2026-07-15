package comment

import (
	"context"
	"errors"
	"social_network/internal/stats"

	"github.com/google/uuid"
)

var (
	ErrCommentNotFound = errors.New("comment not found")
	ErrForbiddenDelete = errors.New("forbidden: cannot delete someone else's comment")
)

type service struct {
	repo        Repository
	statService stats.StatsService
}

func NewService(repo Repository, statService ...stats.StatsService) (CommentService, error) {
	if repo == nil {
		return nil, ErrRepoEmpty
	}

	return &service{repo: repo, statService: stats.ServiceOrNoop(statService...)}, nil
}

func (s *service) Create(ctx context.Context, userID, postID uuid.UUID, content string) error {
	comment := &Comment{
		ID:      uuid.New(),
		UserID:  userID,
		PostID:  postID,
		Content: content,
	}

	if err := s.statService.IncrementPostComments(ctx, postID); err != nil {
		return err
	}

	return s.repo.Create(ctx, comment)
}

func (s *service) GetByID(ctx context.Context, commentID uuid.UUID) (*Comment, error) {
	return s.repo.GetByID(ctx, commentID)
}

func (s *service) ListByPostID(ctx context.Context, postID uuid.UUID) ([]Comment, error) {
	return s.repo.ListByPostID(ctx, postID)
}

func (s *service) Delete(ctx context.Context, userID, commentID uuid.UUID) error {
	comment, err := s.repo.GetByID(ctx, commentID)
	if err != nil {
		return err
	}

	if err := s.statService.DecrementPostComments(ctx, comment.PostID); err != nil {
		return err
	}
	return s.repo.Delete(ctx, userID, commentID)
}
