package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"social_network/internal/user"
	"social_network/internal/walk"
	"social_network/transport/http/middleware"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

type walkStub struct{ current *walk.Result }

func (s *walkStub) Start(_ context.Context, id uuid.UUID, lat, lon float64, visibility user.LocationVisibility, _ time.Time) (*walk.Result, error) {
	s.current = &walk.Result{ID: uuid.New(), UserID: id, Latitude: lat, Longitude: lon, Visibility: visibility}
	return s.current, nil
}
func (s *walkStub) Current(_ context.Context, _ uuid.UUID, _ time.Time) (*walk.Result, error) {
	if s.current == nil {
		return nil, walk.ErrNotFound
	}
	return s.current, nil
}
func (s *walkStub) UpdateLocation(_ context.Context, _ uuid.UUID, _, _ float64, _ time.Time) (*walk.Result, error) {
	return nil, walk.ErrNotFound
}
func (s *walkStub) UpdateVisibility(_ context.Context, _ uuid.UUID, _ user.LocationVisibility, _ time.Time) (*walk.Result, error) {
	return nil, walk.ErrNotFound
}
func (s *walkStub) Stop(_ context.Context, _ uuid.UUID, _ time.Time) error {
	s.current = nil
	return nil
}
func (s *walkStub) Nearby(_ context.Context, _ uuid.UUID, _, _, _ float64, _ time.Time) ([]walk.Result, error) {
	return []walk.Result{}, nil
}

func TestWalkHandlerRequiresAuthenticationAndUsesJWTIdentity(t *testing.T) {
	stub := &walkStub{}
	handler := NewWalkHandler(walk.NewService(stub))
	body := []byte(`{"latitude":55.75,"longitude":37.61,"visibility":0,"user_id":"` + uuid.NewString() + `"}`)
	request := httptest.NewRequest(http.MethodPost, "/walks/start", bytes.NewReader(body))
	recorder := httptest.NewRecorder()
	handler.Start(recorder, request)
	require.Equal(t, http.StatusUnauthorized, recorder.Code)

	authenticatedID := uuid.New()
	request = httptest.NewRequest(http.MethodPost, "/walks/start", bytes.NewReader(body))
	request = request.WithContext(context.WithValue(request.Context(), middleware.UserIDKey, authenticatedID))
	recorder = httptest.NewRecorder()
	handler.Start(recorder, request)
	require.Equal(t, http.StatusCreated, recorder.Code)
	var result walk.Result
	require.NoError(t, json.Unmarshal(recorder.Body.Bytes(), &result))
	require.Equal(t, authenticatedID, result.UserID)
}

func TestWalkHandlerRejectsMissingCoordinates(t *testing.T) {
	handler := NewWalkHandler(walk.NewService(&walkStub{}))
	request := httptest.NewRequest(http.MethodPost, "/walks/start", bytes.NewBufferString(`{"longitude":37.61,"visibility":0}`))
	request = request.WithContext(context.WithValue(request.Context(), middleware.UserIDKey, uuid.New()))
	recorder := httptest.NewRecorder()
	handler.Start(recorder, request)
	require.Equal(t, http.StatusBadRequest, recorder.Code)
}
