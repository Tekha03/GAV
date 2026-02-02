package auth

import (
	"encoding/json"
	"net/http"

	"gav/dbserver"
	"gav/internal/transport/response"
)

type AuthHandler struct {
	authService Service
}

func NewAuthHandler(authService Service) *AuthHandler {
	return &AuthHandler{authService: authService}
}

func (ah *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var dto dbserver.RegisterDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		response.Error(w, http.StatusBadRequest, err)
		return
	}

	token, err := ah.authService.Register(
		r.Context(),
		dto.Email,
		dto.Password,
	)

	if err != nil {
		response.Error(w, http.StatusBadRequest, err)
		return
	}

	response.JSON(w, http.StatusOK, map[string]string{
		"token": token,
	})
}

func (ah *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var dto dbserver.LoginDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		response.Error(w, http.StatusBadRequest, err)
		return
	}

	token, err := ah.authService.Login(
		r.Context(),
		dto.Email,
		dto.Password,
	)

	if err != nil {
		response.Error(w, http.StatusBadRequest, err)
		return
	}

	response.JSON(w, http.StatusOK, map[string]string{
		"token": token,
	})
}
