package handlers

import (
	"encoding/json"
	"net/http"

	"gav/internal/auth"
	"gav/internal/transport/response"
)

type AuthHandler struct {
	service auth.Service
}

func NewAuthHandler(service auth.Service) *AuthHandler {
	return &AuthHandler{service: service}
}

type credentials struct {
	Email		string	`json:"email"`
	Password	string	`json:"password"`
}

func (ah *AuthHandler) Register(writer http.ResponseWriter, reader *http.Request) {
	var req credentials
	if err := json.NewDecoder(reader.Body).Decode(&req); err != nil {
		response.Error(writer, http.StatusBadRequest, err)
		return
	}

	token, err := ah.service.Register(reader.Context(), req.Email, req.Password)
	if err != nil {
		response.Error(writer, http.StatusBadRequest, err)
		return
	}

	response.JSON(writer, http.StatusCreated, map[string]string{"token": token})
}

func (ah *AuthHandler) Login(writer http.ResponseWriter, reader *http.Request) {
	var req credentials
	if err := json.NewDecoder(reader.Body).Decode(&req); err != nil {
		response.Error(writer, http.StatusBadRequest, err)
		return
	}

	token, err := ah.service.Login(reader.Context(), req.Email, req.Password)
	if err != nil {
		response.Error(writer, http.StatusUnauthorized, err)
		return
	}

	response.JSON(writer, http.StatusOK, map[string]string{"token": token})
}
