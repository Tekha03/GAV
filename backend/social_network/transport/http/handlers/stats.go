// read-only
package handlers

import (
	"net/http"
	"social_network/internal/stats"

	"social_network/transport/response"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type StatsHandler struct {
	service stats.StatsService
}

func NewStatsHandler(service stats.StatsService) (*StatsHandler, error) {
	if service == nil {
		return nil, ErrStatsNil
	}

	return &StatsHandler{service: service}, nil
}

// UserStats возвращает статистику пользователя по его ID
// @Summary Получить статистику пользователя
// @Description Возвращает информацию о количестве постов, подписчиков, подписок и др. для конкретного пользователя
// @Tags stats
// @Param userID path string true "UUID пользователя"
// @Produce json
// @Success 200 {object} stats.UserStats
// @Failure 400 {object} response.ErrorResponse "Неверный UUID"
// @Failure 500 {object} response.ErrorResponse "Внутренняя ошибка сервиса"
// @Router /stats/users/{userID} [get]
func (h *StatsHandler) UserStats(w http.ResponseWriter, r *http.Request) {
	userID, err := uuid.Parse(chi.URLParam(r, "userID"))
	if err != nil {
		response.Error(w, ErrInvalidInput)
		return
	}

	userStats, err := h.service.UserStats(r.Context(), userID)
	if err != nil {
		response.Error(w, err)
		return
	}

	response.JSON(w, http.StatusOK, userStats)
}

// UserStats возвращает статистику пользователя по его ID
// @Summary Получить статистику пользователя
// @Description Возвращает информацию о количестве постов, подписчиков, подписок и др. для конкретного пользователя
// @Tags stats
// @Param userID path string true "UUID пользователя"
// @Produce json
// @Success 200 {object} stats.UserStats
// @Failure 400 {object} response.ErrorResponse "Неверный UUID"
// @Failure 500 {object} response.ErrorResponse "Внутренняя ошибка сервиса"
// @Router /stats/users/{userID} [get]
func (h *StatsHandler) ProfileStats(w http.ResponseWriter, r *http.Request) {
	userID, err := uuid.Parse(chi.URLParam(r, "userID"))
	if err != nil {
		response.Error(w, ErrInvalidInput)
		return
	}

	profileStats, err := h.service.ProfileStats(r.Context(), userID)
	if err != nil {
		response.Error(w, err)
		return
	}

	response.JSON(w, http.StatusOK, profileStats)
}

// UserStats возвращает статистику пользователя по его ID
// @Summary Получить статистику пользователя
// @Description Возвращает информацию о количестве постов, подписчиков, подписок и др. для конкретного пользователя
// @Tags stats
// @Param userID path string true "UUID пользователя"
// @Produce json
// @Success 200 {object} stats.UserStats
// @Failure 400 {object} response.ErrorResponse "Неверный UUID"
// @Failure 500 {object} response.ErrorResponse "Внутренняя ошибка сервиса"
// @Router /stats/users/{userID} [get]
func (h *StatsHandler) PostStats(w http.ResponseWriter, r *http.Request) {
	postID, err := uuid.Parse(chi.URLParam(r, "postID"))
	if err != nil {
		response.Error(w, ErrInvalidInput)
		return
	}

	postStats, err := h.service.PostStats(r.Context(), postID)
	if err != nil {
		response.Error(w, err)
		return
	}

	response.JSON(w, http.StatusOK, postStats)
}
