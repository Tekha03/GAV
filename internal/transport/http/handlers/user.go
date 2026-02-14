package handlers

import (
	"net/http"

	"gav/internal/transport/http/dto"
	"gav/internal/transport/http/middleware"
	"gav/internal/transport/response"
	"gav/internal/user"
)

type UserHandler struct {
	service user.UserService
}

func NewUserHandler(service user.UserService) *UserHandler {
	return &UserHandler{service: service}
}

func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserID(r.Context())
	if !ok {
		response.Error(w, ErrUnauthorized)
		return
	}

	user, err := h.service.GetByID(r.Context(), userID)
	if err != nil {
		response.Error(w, err)
		return
	}

	dtoUser := dto.NewUserResponse(user)
	response.JSON(w, http.StatusOK, dtoUser)
}
