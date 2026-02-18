// read-only
package handlers

import (
	"net/http"

	"gav/internal/stats"
	"gav/transport/http/middleware"
	"gav/transport/response"
)

type StatsHandler struct {
	service stats.StatsService
}

func NewStatsHandler(service stats.StatsService) *StatsHandler {
	return &StatsHandler{service: service}
}

func (h *StatsHandler) Get(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserID(r.Context())
	if !ok {
		response.Error(w, ErrUnauthorized)
		return
	}

	stats, err := h.service.Get(r.Context(), userID)
	if err != nil {
		response.Error(w, err)
		return
	}

	response.JSON(w, http.StatusOK, stats)
}
