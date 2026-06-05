package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"social_network/internal/settings"
	"social_network/transport/http/middleware"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func setupSettingsHandler(t *testing.T) (*SettingsHandler, *MockSettingsServcie) {
	service := new(MockSettingsServcie)

	handler, err := NewSettingsHandler(service)
	require.NoError(t, err)

	return handler, service
}

func withUserSettings(ctx context.Context, userID uuid.UUID) context.Context {
	return context.WithValue(ctx, middleware.UserIDKey, userID)
}

func TestSettingsHandler_Get_Success(t *testing.T) {
	handler, service := setupSettingsHandler(t)

	userID := uuid.New()

	expected := &settings.UserSettings{
		UserID:        userID,
		ProfilePrivacy: true,
		ShowLocation:  false,
		AllowMessages: true,
	}

	service.
		On("Get", mock.Anything, userID).
		Return(expected, nil)

	req := httptest.NewRequest(http.MethodGet, "/settings", nil)
	req = req.WithContext(withUser(req.Context(), userID))

	w := httptest.NewRecorder()

	handler.Get(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	service.AssertExpectations(t)
}

func TestSettingsHandler_Get_Unauthorized(t *testing.T) {
	handler, _ := setupSettingsHandler(t)

	req := httptest.NewRequest(http.MethodGet, "/settings", nil)
	w := httptest.NewRecorder()

	handler.Get(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestSettingsHandler_Update_Success(t *testing.T) {
	handler, service := setupSettingsHandler(t)

	userID := uuid.New()

	profilePrivacy := true
	showLocation := true
	allowMessages := false

	input := settings.UpdateSettingsInput{
		ProfilePrivacy:  &profilePrivacy,
		ShowLocation:   &showLocation,
		AllowMessages:  &allowMessages,
	}

	service.
		On("Update", mock.Anything, userID, input).
		Return(nil)

	body, _ := json.Marshal(input)

	req := httptest.NewRequest(http.MethodPut, "/settings", bytes.NewReader(body))
	req = req.WithContext(withUser(req.Context(), userID))

	w := httptest.NewRecorder()

	handler.Update(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)
	service.AssertExpectations(t)
}

func TestSettingsHandler_Update_Unauthorized(t *testing.T) {
	handler, _ := setupSettingsHandler(t)

	req := httptest.NewRequest(http.MethodPut, "/settings", nil)
	w := httptest.NewRecorder()

	handler.Update(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestSettingsHandler_Update_InvalidJSON(t *testing.T) {
	handler, _ := setupSettingsHandler(t)

	userID := uuid.New()

	req := httptest.NewRequest(http.MethodPut, "/settings", bytes.NewBuffer([]byte("{invalid")))
	req = req.WithContext(withUser(req.Context(), userID))

	w := httptest.NewRecorder()

	handler.Update(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}
