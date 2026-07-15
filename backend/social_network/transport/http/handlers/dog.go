package handlers

import (
	"encoding/json"
	"net/http"
	"social_network/internal/dog"
	"social_network/internal/validation"
	"social_network/transport/http/middleware"
	"social_network/transport/response"

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

func (h *DogHandler) ListMine(w http.ResponseWriter, r *http.Request) {
	ownerID, ok := middleware.UserID(r.Context())
	if !ok {
		response.Error(w, ErrUnauthorized)
		return
	}

	if queryOwnerID := r.URL.Query().Get("owner_id"); queryOwnerID != "" {
		parsedOwnerID, err := uuid.Parse(queryOwnerID)
		if err != nil {
			response.Error(w, ErrInvalidInput)
			return
		}
		ownerID = parsedOwnerID
	} else if queryUserID := r.URL.Query().Get("user_id"); queryUserID != "" {
		parsedUserID, err := uuid.Parse(queryUserID)
		if err != nil {
			response.Error(w, ErrInvalidInput)
			return
		}
		ownerID = parsedUserID
	}

	dogs, err := h.service.ListByOwnerID(r.Context(), ownerID)
	if err != nil {
		response.Error(w, err)
		return
	}

	response.JSON(w, http.StatusOK, dogs)
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
