package handlers

import (
	"errors"
	"net/http"

	"gav/internal/transport/http/middleware"
	"gav/internal/transport/response"
	"gav/internal/user"
)

type UserHandler struct {
	service user.Service
}

func NewUserHandler(service user.Service) *UserHandler {
	return &UserHandler{service: service}
}

func (h *UserHandler) Me(writer http.ResponseWriter, reader *http.Request) {
	userID, ok := middleware.UserID(reader.Context())
	if !ok {
		response.Error(writer, http.StatusUnauthorized, errors.New("unauthorized"))
		return
	}

	user, err := h.service.GetByID(reader.Context(), userID)
	if err != nil {
		response.Error(writer, http.StatusNotFound, err)
		return
	}

	response.JSON(writer, http.StatusOK, user)
}
