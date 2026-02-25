package handlers

import (
	"encoding/json"
	"net/http"

	"gav/internal/settings"
	"gav/internal/validation"
	"gav/transport/http/middleware"
	"gav/transport/response"
)

type SettingsHandler struct {
	service settings.SettingsService
}

func NewSettingsHandler(service settings.SettingsService) (*SettingsHandler, error) {
	if service == nil {
		return nil, ErrSettingsNil
	}

	return &SettingsHandler{service: service}, nil
}

func (h *SettingsHandler) Get(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserID(r.Context())
	if !ok {
		response.Error(w, ErrUnauthorized)
		return
	}

	settings, err := h.service.Get(r.Context(), userID)
	if err != nil {
		response.Error(w, err)
		return
	}

	response.JSON(w, http.StatusOK, settings)
}

func (h *SettingsHandler) Update(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserID(r.Context())
	if !ok {
		response.Error(w, ErrUnauthorized)
		return
	}

	var input settings.UpdateSettingsInput
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
