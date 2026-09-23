package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	apperrors "shared/app_errors"
	"social_network/internal/user"
	"social_network/internal/walk"
	"social_network/transport/http/middleware"
	"social_network/transport/response"
)

type WalkHandler struct{ service *walk.Service }

func NewWalkHandler(service *walk.Service) *WalkHandler { return &WalkHandler{service: service} }

type walkInput struct {
	Latitude   *float64                 `json:"latitude"`
	Longitude  *float64                 `json:"longitude"`
	Visibility *user.LocationVisibility `json:"visibility"`
}

func (h *WalkHandler) Start(w http.ResponseWriter, r *http.Request) {
	id, ok := middleware.UserID(r.Context())
	if !ok {
		response.Error(w, ErrUnauthorized)
		return
	}
	var input walkInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil || input.Latitude == nil || input.Longitude == nil || input.Visibility == nil {
		response.Error(w, ErrInvalidInput)
		return
	}
	result, err := h.service.Start(r.Context(), id, *input.Latitude, *input.Longitude, *input.Visibility)
	if err != nil {
		walkError(w, err)
		return
	}
	response.JSON(w, http.StatusCreated, result)
}

func (h *WalkHandler) Current(w http.ResponseWriter, r *http.Request) {
	id, ok := middleware.UserID(r.Context())
	if !ok {
		response.Error(w, ErrUnauthorized)
		return
	}
	result, err := h.service.Current(r.Context(), id)
	if err != nil {
		walkError(w, err)
		return
	}
	response.JSON(w, http.StatusOK, result)
}

func (h *WalkHandler) UpdateLocation(w http.ResponseWriter, r *http.Request) {
	id, ok := middleware.UserID(r.Context())
	if !ok {
		response.Error(w, ErrUnauthorized)
		return
	}
	var input walkInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil || input.Latitude == nil || input.Longitude == nil {
		response.Error(w, ErrInvalidInput)
		return
	}
	result, err := h.service.UpdateLocation(r.Context(), id, *input.Latitude, *input.Longitude)
	if err != nil {
		walkError(w, err)
		return
	}
	response.JSON(w, http.StatusOK, result)
}

func (h *WalkHandler) UpdateVisibility(w http.ResponseWriter, r *http.Request) {
	id, ok := middleware.UserID(r.Context())
	if !ok {
		response.Error(w, ErrUnauthorized)
		return
	}
	var input walkInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil || input.Visibility == nil {
		response.Error(w, ErrInvalidInput)
		return
	}
	result, err := h.service.UpdateVisibility(r.Context(), id, *input.Visibility)
	if err != nil {
		walkError(w, err)
		return
	}
	response.JSON(w, http.StatusOK, result)
}

func (h *WalkHandler) Stop(w http.ResponseWriter, r *http.Request) {
	id, ok := middleware.UserID(r.Context())
	if !ok {
		response.Error(w, ErrUnauthorized)
		return
	}
	if err := h.service.Stop(r.Context(), id); err != nil {
		walkError(w, err)
		return
	}
	response.JSON(w, http.StatusNoContent, nil)
}

func (h *WalkHandler) Nearby(w http.ResponseWriter, r *http.Request) {
	results, ok := h.nearby(w, r)
	if !ok {
		return
	}
	response.JSON(w, http.StatusOK, results)
}

// NearbyDogs keeps the older endpoint private by using the same SQL visibility filter.
func (h *WalkHandler) NearbyDogs(w http.ResponseWriter, r *http.Request) {
	results, ok := h.nearby(w, r)
	if !ok {
		return
	}
	type dogResult struct {
		walk.Dog
		OwnerID        string  `json:"owner_id"`
		Lat            float64 `json:"lat"`
		Lon            float64 `json:"lon"`
		DistanceMeters float64 `json:"distance_meters"`
	}
	dogs := make([]dogResult, 0)
	for _, result := range results {
		for _, dog := range result.Dogs {
			dogs = append(dogs, dogResult{Dog: dog, OwnerID: result.UserID.String(), Lat: result.Latitude, Lon: result.Longitude, DistanceMeters: result.DistanceMeters})
		}
	}
	response.JSON(w, http.StatusOK, dogs)
}

func (h *WalkHandler) nearby(w http.ResponseWriter, r *http.Request) ([]walk.Result, bool) {
	id, ok := middleware.UserID(r.Context())
	if !ok {
		response.Error(w, ErrUnauthorized)
		return nil, false
	}
	query := r.URL.Query()
	lat, errLat := strconv.ParseFloat(query.Get("lat"), 64)
	lon, errLon := strconv.ParseFloat(query.Get("lon"), 64)
	radius, errRadius := strconv.ParseFloat(query.Get("radius"), 64)
	if errLat != nil || errLon != nil || errRadius != nil {
		response.Error(w, ErrInvalidInput)
		return nil, false
	}
	results, err := h.service.Nearby(r.Context(), id, lat, lon, radius)
	if err != nil {
		walkError(w, err)
		return nil, false
	}
	return results, true
}

func walkError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, walk.ErrInvalidCoordinates), errors.Is(err, walk.ErrInvalidVisibility), errors.Is(err, walk.ErrInvalidRadius):
		response.Error(w, apperrors.New(apperrors.Validation, err.Error()))
	case errors.Is(err, walk.ErrNotFound):
		response.Error(w, apperrors.New(apperrors.WalkNotFound, err.Error()))
	case errors.Is(err, user.ErrUserNotFound):
		response.Error(w, err)
	case errors.Is(err, walk.ErrAlreadyActive):
		response.Error(w, apperrors.New(apperrors.WalkAlreadyActive, err.Error()))
	default:
		response.Error(w, err)
	}
}
