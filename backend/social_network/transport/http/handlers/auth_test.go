package handlers

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"social_network/internal/auth"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestAuthHandler_Register(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockSvc := new(MockAuthService)
		h, _ := NewAuthHandler(mockSvc)

		body := `{"email":"test@mail.com","password":"123456"}`

		mockSvc.On("Register", mock.Anything, "test@mail.com", "123456").
			Return(&auth.AuthTokens{AccessToken: "token123"}, nil)

		req := httptest.NewRequest(http.MethodPost, "/auth/register", bytes.NewBufferString(body))
		w := httptest.NewRecorder()

		h.Register(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
		assert.Contains(t, w.Body.String(), "token123")
	})

	t.Run("invalid json", func(t *testing.T) {
		mockSvc := new(MockAuthService)
		h, _ := NewAuthHandler(mockSvc)

		req := httptest.NewRequest(http.MethodPost, "/auth/register", bytes.NewBufferString("{invalid"))
		w := httptest.NewRecorder()

		h.Register(w, req)

		assert.NotEqual(t, http.StatusCreated, w.Code)
	})

	t.Run("service error", func(t *testing.T) {
		mockSvc := new(MockAuthService)
		h, _ := NewAuthHandler(mockSvc)

		body := `{"email":"test@mail.com","password":"123456"}`

		mockSvc.On("Register", mock.Anything, "test@mail.com", "123456").
			Return(&auth.AuthTokens{}, ErrFail)

		req := httptest.NewRequest(http.MethodPost, "/auth/register", bytes.NewBufferString(body))
		w := httptest.NewRecorder()

		h.Register(w, req)

		assert.NotEqual(t, http.StatusCreated, w.Code)
	})
}

func TestAuthHandler_Refresh(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockSvc := new(MockAuthService)
		h, _ := NewAuthHandler(mockSvc)

		body := `{"refresh_token":"abc"}`

		mockSvc.On("Refresh", mock.Anything, "abc").
			Return(&auth.AuthTokens{AccessToken: "newtoken"}, nil)

		req := httptest.NewRequest(http.MethodPost, "/auth/refresh", bytes.NewBufferString(body))
		w := httptest.NewRecorder()

		h.Refresh(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "newtoken")
	})

	t.Run("invalid json", func(t *testing.T) {
		mockSvc := new(MockAuthService)
		h, _ := NewAuthHandler(mockSvc)

		req := httptest.NewRequest(http.MethodPost, "/auth/refresh", bytes.NewBufferString("{bad"))
		w := httptest.NewRecorder()

		h.Refresh(w, req)

		assert.NotEqual(t, http.StatusOK, w.Code)
	})
}

func TestAuthHandler_Logout(t *testing.T) {
	mockSvc := new(MockAuthService)
	h, _ := NewAuthHandler(mockSvc)

	req := httptest.NewRequest(http.MethodPost, "/auth/logout", nil)
	w := httptest.NewRecorder()

	h.Logout(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "logged out")
}
