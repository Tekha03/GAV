package follow

import (
	"encoding/json"
	"errors"
	"net/http"

	"gav/internal/transport/http/middleware"
	"gav/internal/transport/response"
)

type FollowHandler struct {
	service FollowService
}

func NewFollowHandler(service FollowService) *FollowHandler {
	return &FollowHandler{service: service}
}

func (fh *FollowHandler) Follow(w http.ResponseWriter, r *http.Request) {
	var request struct {
		UserID	uint	`json:"user_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		response.Error(w, http.StatusBadRequest, errors.New("bad request"))
		return
	}

	userID, _ := middleware.UserID(r.Context())
	if err := fh.service.Follow(r.Context(), userID, request.UserID); err != nil {
		response.Error(w, http.StatusBadRequest, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
