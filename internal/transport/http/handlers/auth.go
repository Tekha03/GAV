package handlers

import (
	"encoding/json"
	"net/http"

	"gav/internal/auth"
	"gav/internal/transport/http/dto"
	"gav/internal/transport/response"
	"gav/internal/validation"
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

		response.Error(writer, err)
		return
	}

	if err := validation.Validate(&req); err != nil {
		response.Error(writer, err)
		return
	}

	token, err := ah.service.Register(reader.Context(), req.Email, req.Password)
	if err != nil {
		response.Error(writer, err)
		return
	}

	dtoAuth := dto.AuthResponse{Token: token}
	response.JSON(writer, http.StatusCreated, dtoAuth)
}

func (ah *AuthHandler) Login(writer http.ResponseWriter, reader *http.Request) {
	var req credentials
	if err := json.NewDecoder(reader.Body).Decode(&req); err != nil {
		response.Error(writer, err)
		return
	}

	if err := validation.Validate(&req); err != nil {
		response.Error(writer, err)
		return
	}

	token, err := ah.service.Login(reader.Context(), req.Email, req.Password)
	if err != nil {
		response.Error(writer, err)
		return
	}

	dtoAuth := dto.AuthResponse{Token: token}
	response.JSON(writer, http.StatusOK, dtoAuth)
}
