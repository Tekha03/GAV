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

// Create godoc
// @Summary Create a vaccination record for a dog
// @Description Создать запись о вакцинации для конкретной собаки
// @Tags vaccination
// @Accept json
// @Produce json
// @Param id path string true "Dog ID"
// @Param vaccination body vaccination.CreateVaccinationInput true "Vaccination input"
// @Success 201 {object} vaccination.Vaccination
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /dogs/{id}/vaccinations [post]
func (h *VaccinationHandler) Create(w http.ResponseWriter, r *http.Request) {
	dogID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		response.Error(w, ErrInvalidInput)
		return
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

// Create godoc
// @Summary Create a vaccination record for a dog
// @Description Создать запись о вакцинации для конкретной собаки
// @Tags vaccination
// @Accept json
// @Produce json
// @Param id path string true "Dog ID"
// @Param vaccination body vaccination.CreateVaccinationInput true "Vaccination input"
// @Success 201 {object} vaccination.Vaccination
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /dogs/{id}/vaccinations [post]
func (h *VaccinationHandler) ListByDogID(w http.ResponseWriter, r *http.Request) {
	dogID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		response.Error(w, ErrInvalidInput)
		return
	}

	vaccination, err := h.service.ListByDogID(r.Context(), dogID)
	if err != nil {
		response.Error(w, err)
		return
	}

	response.JSON(w, http.StatusOK, vaccination)
}

// Update godoc
// @Summary Update a vaccination record
// @Description Обновить запись о вакцинации
// @Tags vaccination
// @Accept json
// @Produce json
// @Param id path string true "Dog ID"
// @Param vaccinationID path string true "Vaccination ID"
// @Param vaccination body vaccination.UpdateVaccinationInput true "Vaccination update input"
// @Success 200
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /dogs/{id}/vaccinations/{vaccinationID} [put]
func (h *VaccinationHandler) Update(w http.ResponseWriter, r *http.Request) {
	dogID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		response.Error(w, ErrInvalidInput)
		return
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

// Delete godoc
// @Summary Delete a vaccination record
// @Description Удалить запись о вакцинации
// @Tags vaccination
// @Produce json
// @Param id path string true "Vaccination ID"
// @Success 204
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /vaccinations/{id} [delete]
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
