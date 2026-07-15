package handlers

import (
	"encoding/json"
	"net/http"
	"social_network/internal/settings"

	"social_network/internal/validation"
	"social_network/transport/http/middleware"
	"social_network/transport/response"
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

// Get получает настройки текущего пользователя
// @Summary Получить настройки пользователя
// @Description Возвращает текущие настройки пользователя (ProfilePrivacy, ShowLocation, AllowMessages)
// @Tags settings
// @Security ApiKeyAuth
// @Produce json
// @Success 200 {object} settings.UserSettings
// @Failure 401 {object} response.ErrorResponse "Unauthorized"
// @Failure 500 {object} response.ErrorResponse "Internal Server Error"
// @Router /settings [get]
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

// Update обновляет настройки текущего пользователя
// @Summary Обновить настройки пользователя
// @Description Обновляет настройки пользователя (ProfilePrivacy, ShowLocation, AllowMessages)
// @Tags settings
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param settings body settings.UpdateSettingsInput true "Новые настройки пользователя"
// @Success 204 "No Content"
// @Failure 400 {object} response.ErrorResponse "Invalid input"
// @Failure 401 {object} response.ErrorResponse "Unauthorized"
// @Failure 500 {object} response.ErrorResponse "Internal Server Error"
// @Router /settings [put]
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
