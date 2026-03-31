package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"social_network/internal/dog"
	"social_network/transport/http/middleware"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestDogHandler_Create(t *testing.T) {
	userID := uuid.New()

	t.Run("success", func(t *testing.T) {
		mockSvc := new(MockDogService)
		h, _ := NewDogHandler(mockSvc)

		input := dog.CreateDogInput{Name: "Bob"}
		mockSvc.On("Create", mock.Anything, userID, input).
			Return(&dog.Dog{ID: uuid.New(), Name: "Bob"}, nil)

		body, _ := json.Marshal(input)
		req := httptest.NewRequest(http.MethodPost, "/dogs", bytes.NewBuffer(body))
		req = req.WithContext(context.WithValue(req.Context(), middleware.UserIDKey, userID))

		w := httptest.NewRecorder()
		h.Create(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("unauthorized", func(t *testing.T) {
		mockSvc := new(MockDogService)
		h, _ := NewDogHandler(mockSvc)

		req := httptest.NewRequest(http.MethodPost, "/dogs", nil)
		w := httptest.NewRecorder()

		h.Create(w, req)
		assert.NotEqual(t, http.StatusOK, w.Code)
	})

	t.Run("invalid json", func(t *testing.T) {
		mockSvc := new(MockDogService)
		h, _ := NewDogHandler(mockSvc)

		req := httptest.NewRequest(http.MethodPost, "/dogs", bytes.NewBufferString("{bad"))
		req = req.WithContext(context.WithValue(req.Context(), middleware.UserIDKey, userID))

		w := httptest.NewRecorder()
		h.Create(w, req)

		assert.NotEqual(t, http.StatusOK, w.Code)
	})

	t.Run("service error", func(t *testing.T) {
		mockSvc := new(MockDogService)
		h, _ := NewDogHandler(mockSvc)

		input := dog.CreateDogInput{Name: "Bob"}
		mockSvc.On("Create", mock.Anything, userID, input).Return(&dog.Dog{}, ErrFail)

		body, _ := json.Marshal(input)
		req := httptest.NewRequest(http.MethodPost, "/dogs", bytes.NewBuffer(body))
		req = req.WithContext(context.WithValue(req.Context(), middleware.UserIDKey, userID))

		w := httptest.NewRecorder()
		h.Create(w, req)

		assert.NotEqual(t, http.StatusOK, w.Code)
	})
}

func TestDogHandler_GetPublic(t *testing.T) {
	dogID := uuid.New()

	t.Run("success", func(t *testing.T) {
		mockSvc := new(MockDogService)
		h, _ := NewDogHandler(mockSvc)

		mockSvc.On("GetPublic", mock.Anything, dogID).
			Return(&dog.Dog{ID: dogID, Name: "Bob"}, nil)

		r := chi.NewRouter()
		r.Get("/dogs/{id}", h.GetPublic)

		req := httptest.NewRequest(http.MethodGet, "/dogs/"+dogID.String(), nil)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("invalid id", func(t *testing.T) {
		mockSvc := new(MockDogService)
		h, _ := NewDogHandler(mockSvc)

		r := chi.NewRouter()
		r.Get("/dogs/{id}", h.GetPublic)

		req := httptest.NewRequest(http.MethodGet, "/dogs/invalid", nil)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)
		assert.NotEqual(t, http.StatusOK, w.Code)
	})
}

func TestDogHandler_GetPrivate(t *testing.T) {
	userID := uuid.New()
	dogID := uuid.New()

	t.Run("success", func(t *testing.T) {
		mockSvc := new(MockDogService)
		h, _ := NewDogHandler(mockSvc)

		mockSvc.On("GetPrivate", mock.Anything, userID, dogID).
			Return(&dog.Dog{ID: dogID}, nil)

		r := chi.NewRouter()
		r.Get("/dogs/{id}/private", h.GetPrivate)

		req := httptest.NewRequest(http.MethodGet, "/dogs/"+dogID.String()+"/private", nil)
		req = req.WithContext(context.WithValue(req.Context(), middleware.UserIDKey, userID))

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("unauthorized", func(t *testing.T) {
		mockSvc := new(MockDogService)
		h, _ := NewDogHandler(mockSvc)

		r := chi.NewRouter()
		r.Get("/dogs/{id}/private", h.GetPrivate)

		req := httptest.NewRequest(http.MethodGet, "/dogs/"+dogID.String()+"/private", nil)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)
		assert.NotEqual(t, http.StatusOK, w.Code)
	})
}

func TestDogHandler_Update(t *testing.T) {
	userID := uuid.New()
	dogID := uuid.New()

	t.Run("success", func(t *testing.T) {
		mockSvc := new(MockDogService)
		h, _ := NewDogHandler(mockSvc)

		dogName := "New"
		input := dog.UpdateDogInput{Name: &dogName}
		mockSvc.On("Update", mock.Anything, userID, dogID, input).Return(nil)

		r := chi.NewRouter()
		r.Put("/dogs/{id}", h.Update)

		body, _ := json.Marshal(input)
		req := httptest.NewRequest(http.MethodPut, "/dogs/"+dogID.String(), bytes.NewBuffer(body))
		req = req.WithContext(context.WithValue(req.Context(), middleware.UserIDKey, userID))

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNoContent, w.Code)
	})

	t.Run("invalid json", func(t *testing.T) {
		mockSvc := new(MockDogService)
		h, _ := NewDogHandler(mockSvc)

		r := chi.NewRouter()
		r.Put("/dogs/{id}", h.Update)

		req := httptest.NewRequest(http.MethodPut, "/dogs/"+dogID.String(), bytes.NewBufferString("{bad"))
		req = req.WithContext(context.WithValue(req.Context(), middleware.UserIDKey, userID))

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.NotEqual(t, http.StatusNoContent, w.Code)
	})
}

func TestDogHandler_Delete(t *testing.T) {
	userID := uuid.New()
	dogID := uuid.New()

	t.Run("success", func(t *testing.T) {
		mockSvc := new(MockDogService)
		h, _ := NewDogHandler(mockSvc)

		mockSvc.On("Delete", mock.Anything, userID, dogID).Return(nil)

		r := chi.NewRouter()
		r.Delete("/dogs/{id}", h.Delete)

		req := httptest.NewRequest(http.MethodDelete, "/dogs/"+dogID.String(), nil)
		req = req.WithContext(context.WithValue(req.Context(), middleware.UserIDKey, userID))

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNoContent, w.Code)
	})

	t.Run("service error", func(t *testing.T) {
		mockSvc := new(MockDogService)
		h, _ := NewDogHandler(mockSvc)

		mockSvc.On("Delete", mock.Anything, userID, dogID).Return(ErrFail)

		r := chi.NewRouter()
		r.Delete("/dogs/{id}", h.Delete)

		req := httptest.NewRequest(http.MethodDelete, "/dogs/"+dogID.String(), nil)
		req = req.WithContext(context.WithValue(req.Context(), middleware.UserIDKey, userID))

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.NotEqual(t, http.StatusNoContent, w.Code)
	})
}

func TestDogHandler_FindDogsNearby_Success(t *testing.T) {
	mockService := new(MockDogService)
	handler, _ := NewDogHandler(mockService)

	userID := uuid.New()
	lat, lon, radius := 55.75, 37.61, 1000.0

	expectedDogs := []*dog.Dog{
		{ID: uuid.New(), Name: "Rex", LocationStatus: 1},
	}

	mockService.On("FindDogsNearby", mock.Anything, userID, lat, lon, radius).
		Return(expectedDogs, nil)

	req := httptest.NewRequest(http.MethodGet, "/dogs/nearby?lat=55.75&lon=37.61&radius=1000", nil)
	ctx := context.WithValue(req.Context(), middleware.UserIDKey, userID)
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()
	handler.FindDogsNearby(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	mockService.AssertExpectations(t)
}

func TestDogHandler_FindDogsNearby_InvalidParams(t *testing.T) {
	mockService := new(MockDogService)
	handler, _ := NewDogHandler(mockService)

	req := httptest.NewRequest(http.MethodGet, "/dogs/nearby?lat=bad&lon=37.61&radius=1000", nil)
	ctx := context.WithValue(req.Context(), middleware.UserIDKey, uuid.New())
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()
	handler.FindDogsNearby(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestDogHandler_FindDogsNearby_Unauthorized(t *testing.T) {
	mockService := new(MockDogService)
	handler, _ := NewDogHandler(mockService)

	req := httptest.NewRequest(http.MethodGet, "/dogs/nearby?lat=55.75&lon=37.61&radius=1000", nil)

	rr := httptest.NewRecorder()
	handler.FindDogsNearby(rr, req)

	assert.Equal(t, http.StatusUnauthorized, rr.Code)
}
