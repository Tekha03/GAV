package handlers

import (
	"encoding/json"
	"net/http"
	"social_network/internal/vaccination"

	"social_network/internal/validation"
	"social_network/transport/response"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type VaccinationHandler struct {
	service vaccination.VaccinationService
}

func NewVaccinationHandler(service vaccination.VaccinationService) (*VaccinationHandler, error) {
	if service == nil {
		return nil, ErrVaccinationNil
	}

	return &VaccinationHandler{service: service}, nil
}

func (h *VaccinationHandler) Create(w http.ResponseWriter, r *http.Request) {
	dogID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		response.Error(w, ErrInvalidInput)
	}

	var input vaccination.CreateVaccinationInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		response.Error(w, ErrInvalidInput)
		return
	}

	if err := validation.Validate(&input); err != nil {
		response.Error(w, err)
		return
	}

	vaccination, err :=  h.service.Create(r.Context(), dogID, input);
	if err != nil {
		response.Error(w, err)
		return
	}

	response.JSON(w, http.StatusCreated, vaccination)
}

func (h *VaccinationHandler) ListByDogID(w http.ResponseWriter, r *http.Request) {
	dogID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		response.Error(w, ErrInvalidInput)
	}

	vaccination, err := h.service.ListByDogID(r.Context(), dogID)
	if err != nil {
		response.Error(w, err)
		return
	}

	response.JSON(w, http.StatusOK, vaccination)
}

func (h *VaccinationHandler) Update(w http.ResponseWriter, r *http.Request) {
	dogID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		response.Error(w, ErrInvalidInput)
	}

	vaccinationID, err := uuid.Parse(chi.URLParam(r, "vaccinationID"))
	if err != nil {
		response.Error(w, ErrInvalidInput)
		return
	}

	var input vaccination.UpdateVaccinationInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		response.Error(w, ErrInvalidInput)
		return
	}

	if err := validation.Validate(&input); err != nil {
		response.Error(w, err)
		return
	}

	if err := h.service.Update(r.Context(), vaccinationID, dogID, input); err != nil {
		response.Error(w, err)
		return
	}

	response.JSON(w, http.StatusOK, nil)
}

func (h *VaccinationHandler) Delete(w http.ResponseWriter, r *http.Request) {
	vaccinationID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		response.Error(w, ErrInvalidInput)
		return
	}

	if err := h.service.Delete(r.Context(), vaccinationID); err != nil {
		response.Error(w, err)
		return
	}

	response.JSON(w, http.StatusNoContent, nil)
}
