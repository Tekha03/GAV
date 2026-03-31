package handlers

import (
	"encoding/json"
	"net/http"
	"social_network/internal/dog"
	"social_network/internal/validation"
	"social_network/transport/http/middleware"
	"social_network/transport/response"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type DogHandler struct {
	service dog.DogService
}

func NewDogHandler(service dog.DogService) (*DogHandler, error) {
	if service == nil {
		return nil, ErrDogNil
	}

	return &DogHandler{service: service}, nil
}

// Create godoc
// @Summary Create dog
// @Description Создать собаку
// @Tags dogs
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param input body dog.CreateDogInput true "Dog data"
// @Success 200 {object} dog.Dog
// @Failure 400 {object} response.ErrorResponse
// @Failure 401 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Security BearerAuth
// @Router /dogs [post]
func (h *DogHandler) Create(w http.ResponseWriter, r *http.Request) {
	ownerID, ok := middleware.UserID(r.Context())
	if !ok {
		response.Error(w, ErrUnauthorized)
		return
	}

	var input dog.CreateDogInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		response.Error(w, ErrInvalidInput)
		return
	}

	if err := validation.Validate(&input); err != nil {
		response.Error(w, err)
		return
	}

	dog, err := h.service.Create(r.Context(), ownerID, input)
	if err != nil {
		response.Error(w, err)
		return
	}

	response.JSON(w, http.StatusOK, dog)
}

// Create godoc
// @Summary Create dog
// @Description Создать собаку
// @Tags dogs
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param input body dog.CreateDogInput true "Dog data"
// @Success 200 {object} dog.Dog
// @Failure 400 {object} response.ErrorResponse
// @Failure 401 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Security BearerAuth
// @Router /dogs [post]
func (h *DogHandler) GetPublic(w http.ResponseWriter, r *http.Request) {
	dogID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		response.Error(w, ErrInvalidInput)
		return
	}

	dog, err := h.service.GetPublic(r.Context(), dogID)
	if err != nil {
		response.Error(w, err)
		return
	}

	response.JSON(w, http.StatusOK, dog)
}

// GetPrivate godoc
// @Summary Get private dog
// @Description Получить приватную информацию о собаке (только владелец)
// @Tags dogs
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param id path string true "Dog ID"
// @Success 200 {object} dog.Dog
// @Failure 400 {object} response.ErrorResponse
// @Failure 401 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Security BearerAuth
// @Router /dogs/{id}/private [get]
func (h *DogHandler) GetPrivate(w http.ResponseWriter, r *http.Request) {
	ownerID, ok := middleware.UserID(r.Context())
	if !ok {
		response.Error(w, ErrUnauthorized)
		return
	}

	dogID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		response.Error(w, ErrInvalidInput)
		return
	}

	dog, err := h.service.GetPrivate(r.Context(), ownerID, dogID)
	if err != nil {
		response.Error(w, err)
		return
	}

	response.JSON(w, http.StatusOK, dog)
}

// GetPrivate godoc
// @Summary Get private dog
// @Description Получить приватную информацию о собаке (только владелец)
// @Tags dogs
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param id path string true "Dog ID"
// @Success 200 {object} dog.Dog
// @Failure 400 {object} response.ErrorResponse
// @Failure 401 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Security BearerAuth
// @Router /dogs/{id}/private [get]
func (h *DogHandler) Update(w http.ResponseWriter, r *http.Request) {
	ownerID, ok := middleware.UserID(r.Context())
	if !ok {
		response.Error(w, ErrUnauthorized)
		return
	}

	dogID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		response.Error(w, ErrInvalidInput)
		return
	}

	var input dog.UpdateDogInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		response.Error(w, ErrInvalidInput)
		return
	}

	if err := validation.Validate(&input); err != nil {
		response.Error(w, err)
		return
	}

	if err := h.service.Update(r.Context(), ownerID, dogID, input); err != nil {
		response.Error(w, err)
		return
	}

	response.JSON(w, http.StatusNoContent, nil)
}

// Update godoc
// @Summary Update dog
// @Description Обновить данные собаки
// @Tags dogs
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param id path string true "Dog ID"
// @Param input body dog.UpdateDogInput true "Dog data"
// @Success 204 {object} nil
// @Failure 400 {object} response.ErrorResponse
// @Failure 401 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Security BearerAuth
// @Router /dogs/{id} [put]
func (h *DogHandler) Delete(w http.ResponseWriter, r *http.Request) {
	ownerID, ok := middleware.UserID(r.Context())
	if !ok {
		response.Error(w, ErrUnauthorized)
		return
	}

	dogID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		response.Error(w, ErrInvalidInput)
		return
	}

	if err := h.service.Delete(r.Context(), ownerID, dogID); err != nil {
		response.Error(w, err)
		return
	}

	response.JSON(w, http.StatusNoContent, nil)
}

// FindDogsNearby godoc
// @Summary Find dogs nearby
// @Description Найти собак поблизости в заданном радиусе
// @Tags dogs
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param lat query number true "Широта (Latitude)"
// @Param lon query number true "Долгота (Longitude)"
// @Param radius query number true "Радиус в метрах"
// @Success 200 {array} dog.Dog
// @Failure 400 {object} response.ErrorResponse
// @Failure 401 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Security BearerAuth
// @Router /dogs/nearby [get]
func (h *DogHandler) FindDogsNearby(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserID(r.Context())
	if !ok {
		response.Error(w, ErrUnauthorized)
		return
	}

	q := r.URL.Query()

	lat, errLat := strconv.ParseFloat(q.Get("lat"), 64)
	lon, errLon := strconv.ParseFloat(q.Get("lon"), 64)
	radius, errRad := strconv.ParseFloat(q.Get("radius"), 64)

	if errLat != nil || errLon != nil || errRad != nil {
		response.Error(w, ErrInvalidInput)
		return
	}

	dogs, err := h.service.FindDogsNearby(r.Context(), userID, lat, lon, radius)
	if err != nil {
		response.Error(w, err)
		return
	}

	if dogs == nil {
		dogs = []*dog.Dog{}
	}

	response.JSON(w, http.StatusOK, dogs)
}
