// read-only
package handlers

import (
	"net/http"

	"gav/internal/stats"
	"gav/transport/response"

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
