package handlers

import (
	"gav/internal/post"
	"gav/transport/http/dto"
	"gav/transport/http/middleware"
	"gav/transport/response"
	"net/http"
	"time"
)

const LIMIT_OF_POSTS = 20

type FeedHandler struct {
	service post.PostService
}

func NewFeedHandler(service post.PostService) *FeedHandler {
	return &FeedHandler{service: service}
}

func (h *FeedHandler) GetFeed(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserID(r.Context())
	if !ok {
		response.Error(w, ErrUnauthorized)
	}

	cursorStr := r.URL.Query().Get("cursor")
	limit := LIMIT_OF_POSTS

	var before time.Time
	if cursorStr != "" {
		var err error
		before, err = time.Parse(time.RFC3339Nano, cursorStr)
		if err != nil {
			response.Error(w, err)
			return
		}
	}

	posts, nextCursor, err := h.service.GetFeed(r.Context(), userID, before, limit)
	if err != nil {
		response.Error(w, err)
		return
	}

	resp := dto.FeedResponse{
		Posts:	make([]dto.PostResponse, len(posts)),
		NextCursor:	"",
		HasMore:	nextCursor != time.Time{},
	}

	for i, post := range posts {
		resp.Posts[i] = dto.NewPostResponse(post)
	}

	if resp.HasMore {
		resp.NextCursor = nextCursor.UTC().Format(time.RFC3339Nano)
	}

	response.JSON(w, http.StatusOK, resp)
}
