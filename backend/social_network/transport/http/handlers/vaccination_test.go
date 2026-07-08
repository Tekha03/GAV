package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"social_network/internal/vaccination"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func toJSON(t *testing.T, v interface{}) *bytes.Reader {
	b, err := json.Marshal(v)
	if err != nil {
		t.Fatal(err)
	}
	return bytes.NewReader(b)
}

func TestVaccinationHandler_Create_Success(t *testing.T) {
	mockService := new(MockVaccinationService)
	handler, _ := NewVaccinationHandler(mockService)

	dogID := uuid.New()
	input := vaccination.CreateVaccinationInput{
		Name:   "Rabies",
		DoneAt: time.Now().Truncate(0),
		Notes:  "First dose",
	}

	expected := &vaccination.Vaccination{
		ID:     uuid.New(),
		DogID:  dogID,
		Name:   input.Name,
		DoneAt: input.DoneAt,
		Notes:  input.Notes,
	}

	mockService.On("Create", mock.Anything, dogID, input).Return(expected, nil)

	req := httptest.NewRequest(http.MethodPost, "/dogs/"+dogID.String()+"/vaccinations", toJSON(t, input))
	rc := chi.NewRouteContext()
	rc.URLParams.Add("id", dogID.String())
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rc))

	rr := httptest.NewRecorder()
	handler.Create(rr, req)

	assert.Equal(t, http.StatusCreated, rr.Code)
	mockService.AssertExpectations(t)
}

func TestVaccinationHandler_ListByDogID_Success(t *testing.T) {
	mockService := new(MockVaccinationService)
	handler, _ := NewVaccinationHandler(mockService)

	dogID := uuid.New()
	expected := []*vaccination.Vaccination{
		{ID: uuid.New(), DogID: dogID, Name: "Rabies"},
	}

	mockService.On("ListByDogID", mock.Anything, dogID).Return(expected, nil)

	req := httptest.NewRequest(http.MethodGet, "/dogs/"+dogID.String()+"/vaccinations", nil)
	rc := chi.NewRouteContext()
	rc.URLParams.Add("id", dogID.String())
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rc))

	rr := httptest.NewRecorder()
	handler.ListByDogID(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	mockService.AssertExpectations(t)
}

func TestVaccinationHandler_Update_Success(t *testing.T) {
	mockService := new(MockVaccinationService)
	handler, _ := NewVaccinationHandler(mockService)

	dogID := uuid.New()
	vaccID := uuid.New()
	name := "Distemper"
	input := vaccination.UpdateVaccinationInput{Name: &name}

	mockService.On("Update", mock.Anything, vaccID, dogID, input).Return(nil)

	req := httptest.NewRequest(http.MethodPut, "/dogs/"+dogID.String()+"/vaccinations/"+vaccID.String(), toJSON(t, input))
	rc := chi.NewRouteContext()
	rc.URLParams.Add("id", dogID.String())
	rc.URLParams.Add("vaccinationID", vaccID.String())
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rc))

	rr := httptest.NewRecorder()
	handler.Update(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	mockService.AssertExpectations(t)
}

func TestVaccinationHandler_Delete_Success(t *testing.T) {
	mockService := new(MockVaccinationService)
	handler, _ := NewVaccinationHandler(mockService)

	vaccID := uuid.New()
	mockService.On("Delete", mock.Anything, vaccID).Return(nil)

	req := httptest.NewRequest(http.MethodDelete, "/vaccinations/"+vaccID.String(), nil)
	rc := chi.NewRouteContext()
	rc.URLParams.Add("id", vaccID.String())
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rc))

	rr := httptest.NewRecorder()
	handler.Delete(rr, req)

	assert.Equal(t, http.StatusNoContent, rr.Code)
	mockService.AssertExpectations(t)
}

func TestVaccinationHandler_InvalidUUID(t *testing.T) {
	mockService := new(MockVaccinationService)
	handler, _ := NewVaccinationHandler(mockService)

	req := httptest.NewRequest(http.MethodGet, "/dogs/invalid-uuid/vaccinations", nil)
	rc := chi.NewRouteContext()
	rc.URLParams.Add("id", "invalid-uuid")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rc))

	rr := httptest.NewRecorder()
	handler.ListByDogID(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}
