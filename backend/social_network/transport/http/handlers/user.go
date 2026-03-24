package handlers

import (
	"encoding/json"
	"net/http"
	"social_network/internal/user"

	"social_network/internal/validation"
	"social_network/transport/http/dto"
	"social_network/transport/response"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type UserHandler struct {
	service 	user.UserService
}

func NewUserHandler(service user.UserService) (*UserHandler, error) {
	if service == nil {
		return nil, ErrUserNil
	}

	return &UserHandler{service: service}, nil
}

func (h *UserHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")

	userID, err := uuid.Parse(idParam)
	if err != nil {
		response.Error(w, ErrInvalidInput)
		return
	}

	user, err := h.service.GetByID(r.Context(), userID)
	if err != nil {
		response.Error(w, err)
		return
	}

	response.JSON(w, http.StatusOK, dto.UserResponse{ID: userID, Email: user.Email})
}

func (h *UserHandler) Update(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")

	userID, err := uuid.Parse(idParam)
	if err != nil {
		response.Error(w, ErrInvalidInput)
		return
	}

	var input user.UpdateUserInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		response.Error(w, ErrInvalidInput)
		return
	}

	if err := validation.Validate(&input); err != nil {
		response.Error(w, err)
		return
	}

	if err := h.service.Update(r.Context(), userID, input); err != nil {
		response.Error(w, err)
		return
	}

	response.JSON(w, http.StatusNoContent, nil)
}

func (h *UserHandler) Delete(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")

	userID, err := uuid.Parse(idParam)
	if err != nil {
		response.Error(w, ErrInvalidInput)
	}

	if err := h.service.Delete(r.Context(), userID); err != nil {
		response.Error(w, err)
		return
	}

	response.JSON(w, http.StatusNoContent, nil)
}
