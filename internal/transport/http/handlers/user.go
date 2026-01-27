package handlers

import (
	"errors"
	"gav/internal/transport/response"
	"gav/internal/transport/http/middleware"
	"gav/storage"
	"net/http"
)

type UserHandler struct {
	userRepo storage.Repository
}

func NewUserHandler(repo storage.Repository) *UserHandler {
	return &UserHandler{userRepo: repo}
}

func (h *UserHandler) Me(writer http.ResponseWriter, reader *http.Request) {
	userID, ok := middleware.UserID(reader.Context())
	if !ok {
		response.Error(writer, http.StatusUnauthorized, errors.New("unauthorized"))
		return
	}

	user, err := h.userRepo.GetByID(reader.Context(), userID)
	if err != nil {
		response.Error(writer, http.StatusNotFound, err)
		return
	}

	response.JSON(writer, http.StatusOK, user)
}
