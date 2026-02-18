package handlers

import (
	"encoding/json"
	"net/http"

	"gav/internal/dog"
	"gav/internal/validation"
	"gav/transport/http/middleware"
	"gav/transport/response"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type DogHandler struct {
	service dog.DogService
}

func NewDog(service dog.DogService) *DogHandler {
	return &DogHandler{service: service}
}

func (h *DogHandler) Create(w http.ResponseWriter, r *http.Request) {
	ownerID, ok := middleware.UserID(r.Context())
	if !ok {
		response.Error(w, ErrUnauthorized)
		return
	}

	var input dog.CreateDogInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		response.Error(w, ErrInvalidInput)
		return
	}

	if err := validation.Validate(&input); err != nil {
		response.Error(w, err)
		return
	}

	dog, err := h.service.Create(r.Context(), ownerID, input)
	if err != nil {
		response.Error(w, err)
		return
	}

	response.JSON(w, http.StatusOK, dog)
}

func (h *DogHandler) GetPublic(w http.ResponseWriter, r *http.Request) {
	dogID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		response.Error(w, ErrInvalidInput)
		return
	}

	dog, err := h.service.GetPublic(r.Context(), dogID)
	if err != nil {
		response.Error(w, err)
		return
	}

	response.JSON(w, http.StatusOK, dog)
}

func (h *DogHandler) GetPrivate(w http.ResponseWriter, r *http.Request) {
	ownerID, ok := middleware.UserID(r.Context())
	if !ok {
		response.Error(w, ErrUnauthorized)
		return
	}

	dogID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		response.Error(w, ErrInvalidInput)
		return
	}

	dog, err := h.service.GetPrivate(r.Context(), ownerID, dogID)
	if err != nil {
		response.Error(w, err)
		return
	}

	response.JSON(w, http.StatusOK, dog)
}

func (h *DogHandler) Update(w http.ResponseWriter, r *http.Request) {
	ownerID, ok := middleware.UserID(r.Context())
	if !ok {
		response.Error(w, ErrUnauthorized)
		return
	}

	dogID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		response.Error(w, ErrInvalidInput)
		return
	}

	var input dog.UpdateDogInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		response.Error(w, ErrInvalidInput)
		return
	}

	if err := validation.Validate(&input); err != nil {
		response.Error(w, err)
		return
	}

	if err := h.service.Update(r.Context(), ownerID, dogID, input); err != nil {
		response.Error(w, err)
		return
	}

	response.JSON(w, http.StatusNoContent, nil)
}

func (h *DogHandler) Delete(w http.ResponseWriter, r *http.Request) {
	ownerID, ok := middleware.UserID(r.Context())
	if !ok {
		response.Error(w, ErrUnauthorized)
		return
	}

	dogID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		response.Error(w, ErrInvalidInput)
		return
	}

	if err := h.service.Delete(r.Context(), ownerID, dogID); err != nil {
		response.Error(w, err)
		return
	}

	response.JSON(w, http.StatusNoContent, nil)
}
