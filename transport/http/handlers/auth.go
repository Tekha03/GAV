package handlers

import (
	"encoding/json"
	"net/http"

	"gav/internal/auth"
	"gav/internal/validation"
	"gav/transport/http/dto"
	"gav/transport/http/middleware"
	"gav/transport/response"
)

type AuthHandler struct {
	service auth.AuthService
}

func NewAuthHandler(service auth.AuthService) (*AuthHandler, error) {
	if service == nil {
		return nil, ErrAuthNil
	}

	return &AuthHandler{service: service}, nil
}

type credentials struct {
	Email		string	`json:"email"`
	Password	string	`json:"password"`
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req credentials
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, err)
		return
	}

	if err := validation.Validate(&req); err != nil {
		response.Error(w, err)
		return
	}

	token, err := h.service.Register(r.Context(), req.Email, req.Password)
	if err != nil {
		response.Error(w, err)
		return
	}

	response.JSON(w, http.StatusCreated, dto.AuthResponse{Token: token})
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req credentials
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, err)
		return
	}

	if err := validation.Validate(&req); err != nil {
		response.Error(w, err)
		return
	}

	token, err := h.service.Login(r.Context(), req.Email, req.Password)
	if err != nil {
		response.Error(w, err)
		return
	}

	response.JSON(w, http.StatusOK, dto.AuthResponse{Token: token})
}

func (h *AuthHandler) Me(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserID(r.Context())
	if !ok {
		response.Error(w, ErrUnauthorized)
	}

	response.JSON(w, http.StatusOK, map[string]any{
		"user_id": userID,
	})
}

// stateless jwt (пока без blacklist)
func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	response.JSON(w, http.StatusOK, map[string]any{
		"message": "logged out",
	})
}
