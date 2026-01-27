package handlers

import (
	"gav/internal/follow"
	"gav/internal/transport/response"
	"gav/internal/validation"
	"net/http"

	"github.com/gin-gonic/gin"
)

type FollowHandler struct {
	service follow.FollowService
}

func NewFollowHandler(service follow.FollowService) *FollowHandler {
	return &FollowHandler{service: service}
}

func (fh *FollowHandler) Follow(ctx *gin.Context) {
	var request struct {
		UserID	uint	`json:"user_id" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&request); err != nil {
		response.Error(ctx.Writer, http.StatusBadRequest, err)
		return
	}

	if err := validation.Validate(&request); err != nil {
		response.Error(ctx.Writer, http.StatusBadRequest, err)
		return
	}

	authID := ctx.GetUint("user_id")
	err := fh.service.Follow(ctx, authID, request.UserID)
	if err != nil {
		response.Error(ctx.Writer, http.StatusBadRequest, err)
		return
	}

	ctx.Status(http.StatusNoContent)
}
