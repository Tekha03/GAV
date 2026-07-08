package handlers

import (
	"net/http"
	"social_network/internal/feed"
	"social_network/transport/http/dto"
	"social_network/transport/http/middleware"
	"social_network/transport/response"
	"strconv"
	"time"
)

const LIMIT_OF_POSTS = 20

type FeedHandler struct {
	service feed.FeedService
}

func NewFeedHandler(service feed.FeedService) (*FeedHandler, error) {
	if service == nil {
		return nil, ErrFeedNil
	}

	return &FeedHandler{service: service}, nil
}

// GetFeed godoc
// @Summary      Get user feed
// @Description  Returns paginated feed for authorized user
// @Tags         feed
// @Accept       json
// @Produce      json
// @Param        cursor query string false "Pagination cursor (RFC3339Nano)"
// @Param        limit  query int    false "Number of posts (default 20)"
// @Success      200 {object} dto.FeedResponse
// @Failure      400 {object} response.ErrorResponse
// @Failure      401 {object} response.ErrorResponse
// @Failure      500 {object} response.ErrorResponse
// @Router       /feed [get]
func (h *FeedHandler) GetFeed(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserID(r.Context())
	if !ok {
		response.Error(w, ErrUnauthorized)
		return
	}

	cursor := r.URL.Query().Get("cursor")
	limitStr := r.URL.Query().Get("limit")

	limit := LIMIT_OF_POSTS
	if limitStr != "" {
		newLimit, err := strconv.Atoi(limitStr)
		if err != nil || newLimit <= 0 {
			response.Error(w, ErrInvalidLimit)
			return
		}
		limit = newLimit
	}

	var before time.Time
	if cursor != "" {
		var err error
		before, err = time.Parse(time.RFC3339Nano, cursor)
		if err != nil {
			response.Error(w, ErrInvalidCursor)
			return
		}
	}

	posts, nextCursor, err := h.service.GetFeed(r.Context(), userID, before, limit)
	if err != nil {
		response.Error(w, err)
		return
	}

	resp := dto.FeedResponse{
		Posts:      make([]dto.PostResponse, len(posts)),
		NextCursor: "",
		HasMore:    nextCursor != time.Time{},
	}

	for i, post := range posts {
		resp.Posts[i] = dto.PostResponse{
			ID:        post.ID,
			AuthorID:  post.UserID,
			Content:   post.Content,
			ImageUrl:  post.ImageUrl,
			CreatedAt: post.CreatedAt,
		}
	}

	if resp.HasMore {
		resp.NextCursor = nextCursor.UTC().Format(time.RFC3339Nano)
	}

	response.JSON(w, http.StatusOK, resp)
}
