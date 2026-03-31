// package handlers

// import (
// 	"context"
// 	"net/http"
// 	"net/http/httptest"
// 	"social_network/internal/dog"
// 	"social_network/transport/http/middleware"
// 	"testing"

// 	"github.com/google/uuid"
// 	"github.com/stretchr/testify/assert"
// 	"github.com/stretchr/testify/mock"
// )

// func TestDogHandler_FindDogsNearby_Success(t *testing.T) {
// 	mockService := new(MockUserService)
// 	handler, _ := NewUserHandler(mockService)

// 	userID := uuid.New()
// 	lat, lon, radius := 55.75, 37.61, 1000.0

// 	expectedDogs := []*dog.Dog{
// 		{ID: uuid.New(), Name: "Rex", LocationStatus: 1},
// 	}

// 	mockService.On("FindDogsNearby", mock.Anything, userID, lat, lon, radius).
// 		Return(expectedDogs, nil)

// 	req := httptest.NewRequest(http.MethodGet, "/dogs/nearby?lat=55.75&lon=37.61&radius=1000", nil)
// 	ctx := context.WithValue(req.Context(), middleware.UserIDKey, userID)
// 	req = req.WithContext(ctx)

// 	rr := httptest.NewRecorder()
// 	handler.FindDogsNearby(rr, req)

// 	assert.Equal(t, http.StatusOK, rr.Code)
// 	mockService.AssertExpectations(t)
// }

// func TestDogHandler_FindDogsNearby_InvalidParams(t *testing.T) {
// 	mockService := new(MockUserService)
// 	handler, _ := NewUserHandler(mockService)

// 	req := httptest.NewRequest(http.MethodGet, "/dogs/nearby?lat=bad&lon=37.61&radius=1000", nil)
// 	ctx := context.WithValue(req.Context(), middleware.UserIDKey, uuid.New())
// 	req = req.WithContext(ctx)

// 	rr := httptest.NewRecorder()
// 	handler.FindDogsNearby(rr, req)

// 	assert.Equal(t, http.StatusBadRequest, rr.Code)
// }

// func TestDogHandler_FindDogsNearby_Unauthorized(t *testing.T) {
// 	mockService := new(MockUserService)
// 	handler, _ := NewDogHandler(mockService)

// 	req := httptest.NewRequest(http.MethodGet, "/dogs/nearby?lat=55.75&lon=37.61&radius=1000", nil)

// 	rr := httptest.NewRecorder()
// 	handler.FindUserNearby(rr, req)

// 	assert.Equal(t, http.StatusUnauthorized, rr.Code)
// }
