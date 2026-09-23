package handlers

import (
	"encoding/json"
	"io"
	"net/http"

	"social_network/internal/device"
	"social_network/transport/http/middleware"
	"social_network/transport/response"
)

type DeviceHandler struct{ service *device.Service }

func NewDeviceHandler(service *device.Service) *DeviceHandler {
	return &DeviceHandler{service: service}
}

func (h *DeviceHandler) Register(w http.ResponseWriter, r *http.Request) {
	h.handle(w, r, true)
}

func (h *DeviceHandler) Unregister(w http.ResponseWriter, r *http.Request) {
	h.handle(w, r, false)
}

func (h *DeviceHandler) handle(w http.ResponseWriter, r *http.Request, register bool) {
	id, ok := middleware.UserID(r.Context())
	if !ok {
		response.Error(w, middleware.ErrUnauthorized)
		return
	}
	var input struct {
		Token string `json:"token"`
	}
	decoder := json.NewDecoder(http.MaxBytesReader(w, r.Body, 8192))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&input); err != nil {
		response.Error(w, device.ErrInvalidToken)
		return
	}
	var extra any
	if err := decoder.Decode(&extra); err != io.EOF {
		response.Error(w, device.ErrInvalidToken)
		return
	}
	var err error
	if register {
		err = h.service.Register(r.Context(), id, input.Token)
	} else {
		err = h.service.Unregister(r.Context(), id, input.Token)
	}
	if err != nil {
		response.Error(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
