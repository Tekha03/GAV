package handlers

import (
	"encoding/json"
	"net/http"

	"gav/internal/profile"
	"gav/internal/transport/http/middleware"
	"gav/internal/transport/response"
	"gav/internal/validation"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type ProfileHandler struct {
	service profile.ProfileService
}

func NewProfileHandler(service profile.ProfileService) *ProfileHandler {
	return &ProfileHandler{service: service}
}

func (h *ProfileHandler) Create(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserID(r.Context())
	if !ok {
		response.Error(w, ErrUnauthorized)
		return
	}

	var input profile.CreateProfileInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		response.Error(w, ErrInvalidInput)
		return
	}

	if err := validation.Validate(&input); err != nil {
		response.Error(w, err)
		return
	}

	profile, err := h.service.Create(r.Context(), userID, input)
	if err != nil {
		response.Error(w, err)
		return
	}

	response.JSON(w, http.StatusCreated, profile)
}

func (h *ProfileHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	profileID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		response.Error(w, ErrInvalidInput)
		return
	}

	profile, err := h.service.GetByID(r.Context(), profileID)
	if err != nil {
		response.Error(w, err)
		return
	}

	response.JSON(w, http.StatusOK, profile)
}

func (h *ProfileHandler) Update(w http.ResponseWriter, r *http.Request) {
	profileID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		response.Error(w, ErrInvalidInput)
		return
	}

	var input profile.UpdateProfileInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		response.Error(w, ErrInvalidInput)
		return
	}

	if err := validation.Validate(&input); err != nil {
		response.Error(w, err)
		return
	}

	if err := h.service.Update(r.Context(), profileID, input); err != nil {
		response.Error(w, err)
		return
	}

	response.JSON(w, http.StatusNoContent, nil)
}

func (h *ProfileHandler) Delete(w http.ResponseWriter, r *http.Request) {
	profileID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		response.Error(w, ErrInvalidInput)
		return
	}

	if err := h.service.Delete(r.Context(), profileID); err != nil {
		response.Error(w, err)
		return
	}

	response.JSON(w, http.StatusNoContent, nil)
}
