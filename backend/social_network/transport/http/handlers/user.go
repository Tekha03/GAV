package handlers

import (
	"encoding/json"
	"net/http"
	"social_network/internal/dog"
	"social_network/internal/user"
	"strconv"

	"social_network/internal/validation"
	"social_network/transport/http/dto"
	"social_network/transport/http/middleware"
	"social_network/transport/response"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type UserHandler struct {
	service user.UserService
}

func NewUserHandler(service user.UserService) (*UserHandler, error) {
	if service == nil {
		return nil, ErrUserNil
	}

	return &UserHandler{service: service}, nil
}

// GetByID godoc
// @Summary Получить пользователя по ID
// @Description Возвращает информацию о пользователе по его UUID
// @Tags users
// @Produce json
// @Param id path string true "UUID пользователя"
// @Param Authorization header string true "Bearer token"
// @Success 200 {object} dto.UserResponse
// @Failure 400 {object} response.ErrorResponse "Некорректный UUID"
// @Failure 401 {object} response.ErrorResponse "Неавторизован"
// @Failure 404 {object} response.ErrorResponse "Пользователь не найден"
// @Failure 500 {object} response.ErrorResponse
// @Security BearerAuth
// @Router /users/{id} [get]
func (h *UserHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")

	userID, err := uuid.Parse(idParam)
	if err != nil {
		response.Error(w, ErrInvalidInput)
		return
	}

	user, err := h.service.GetByID(r.Context(), userID)
	if err != nil {
		response.Error(w, err)
		return
	}

	response.JSON(w, http.StatusOK, dto.UserResponse{ID: userID, Email: user.Email})
}

// GetByID godoc
// @Summary Получить пользователя по ID
// @Description Возвращает информацию о пользователе по его UUID
// @Tags users
// @Produce json
// @Param id path string true "UUID пользователя"
// @Param Authorization header string true "Bearer token"
// @Success 200 {object} dto.UserResponse
// @Failure 400 {object} response.ErrorResponse "Некорректный UUID"
// @Failure 401 {object} response.ErrorResponse "Неавторизован"
// @Failure 404 {object} response.ErrorResponse "Пользователь не найден"
// @Failure 500 {object} response.ErrorResponse
// @Security BearerAuth
// @Router /users/{id} [get]
func (h *UserHandler) Update(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")

	userID, err := uuid.Parse(idParam)
	if err != nil {
		response.Error(w, ErrInvalidInput)
		return
	}

	var input user.UpdateUserInput
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

// Update godoc
// @Summary Обновить данные пользователя
// @Description Обновляет информацию пользователя по UUID
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "UUID пользователя"
// @Param user body user.UpdateUserInput true "Данные для обновления пользователя"
// @Param Authorization header string true "Bearer token"
// @Success 204 "No Content"
// @Failure 400 {object} response.ErrorResponse "Некорректный UUID или входные данные"
// @Failure 401 {object} response.ErrorResponse "Неавторизован"
// @Failure 500 {object} response.ErrorResponse
// @Security BearerAuth
// @Router /users/{id} [put]
func (h *UserHandler) Delete(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")

	userID, err := uuid.Parse(idParam)
	if err != nil {
		response.Error(w, ErrInvalidInput)
		return
	}

	if err := h.service.Delete(r.Context(), userID); err != nil {
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
func (h *UserHandler) FindDogsNearby(w http.ResponseWriter, r *http.Request) {
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

func (h *UserHandler) UpdateLocation(w http.ResponseWriter, r *http.Request) {
	authUserID, ok := middleware.UserID(r.Context())
	if !ok {
		response.Error(w, ErrUnauthorized)
		return
	}

	idParam := chi.URLParam(r, "id")

	userID, err := uuid.Parse(idParam)
	if err != nil {
		response.Error(w, ErrInvalidInput)
		return
	}

	if userID != authUserID {
		response.Error(w, ErrForbidden)
		return
	}

	var input user.UpdateLocationInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		response.Error(w, ErrInvalidInput)
		return
	}

	if err := validation.Validate(&input); err != nil {
		response.Error(w, err)
		return
	}

	if err := h.service.UpdateLocation(r.Context(), userID, input); err != nil {
		response.Error(w, err)
		return
	}

	response.JSON(w, http.StatusNoContent, nil)
}

func (h *UserHandler) SetLocationVisibility(w http.ResponseWriter, r *http.Request) {
	authUserID, ok := middleware.UserID(r.Context())
	if !ok {
		response.Error(w, ErrUnauthorized)
		return
	}

	idParam := chi.URLParam(r, "id")

	userID, err := uuid.Parse(idParam)
	if err != nil {
		response.Error(w, ErrInvalidInput)
		return
	}

	if userID != authUserID {
		response.Error(w, ErrForbidden)
		return
	}

	var input user.SetLocationVisibilityInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		response.Error(w, ErrInvalidInput)
		return
	}

	if err := validation.Validate(&input); err != nil {
		response.Error(w, err)
		return
	}

	if err := h.service.SetLocationVisibility(r.Context(), userID, input); err != nil {
		response.Error(w, err)
		return
	}

	response.JSON(w, http.StatusNoContent, nil)
}
