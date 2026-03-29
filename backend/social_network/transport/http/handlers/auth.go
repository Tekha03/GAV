package handlers

import (
	"encoding/json"
	"net/http"

	"social_network/internal/auth"
	"social_network/internal/validation"
	"social_network/transport/http/dto"
	"social_network/transport/http/middleware"
	"social_network/transport/response"
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

// @Summary      Регистрация пользователя
// @Description  Создает нового пользователя и возвращает JWT-токен
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        input body credentials true "User credentials"
// @Success      201   {object} dto.AuthResponse
// @Failure      400   {object} response.ErrorResponse
// @Failure      500   {object} response.ErrorResponse
// @Router       /auth/register [post]
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

// @Summary      Вход в систему
// @Description  Аутентификация пользователя и получение JWT-токена
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        input body credentials true "User credentials"
// @Success      200   {object} dto.AuthResponse
// @Failure      400   {object} response.ErrorResponse
// @Failure      401   {object} response.ErrorResponse
// @Failure      500   {object} response.ErrorResponse
// @Router       /auth/login [post]
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

// @Summary      Получить текущего пользователя
// @Description  Возвращает ID авторизованного пользователя
// @Tags         auth
// @Produce      json
// @Security     BearerAuth
// @Success      200   {object} map[string]any
// @Failure      401   {object} response.ErrorResponse
// @Router       /auth/me [get]
func (h *AuthHandler) Me(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserID(r.Context())
	if !ok {
		response.Error(w, ErrUnauthorized)
		return
	}

	userInfo, err := h.service.Me(r.Context(), userID)
	if err != nil {
		response.Error(w, err)
		return
	}

	if userInfo == nil {
		response.Error(w, ErrUnauthorized)
		return
	}

	response.JSON(w, http.StatusOK, userInfo)
}

// @Summary      Получить текущего пользователя
// @Description  Возвращает ID авторизованного пользователя
// @Tags         auth
// @Produce      json
// @Security     BearerAuth
// @Success      200   {object} map[string]any
// @Failure      401   {object} response.ErrorResponse
// @Router       /auth/me [get]
func (h *AuthHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	type refreshRequest struct {
		RefreshToken string	`json:"refresh_token" validate:"required"`
	}

	var req refreshRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, err)
		return
	}

	if err := validation.Validate(&req); err != nil {
		response.Error(w, err)
		return
	}

	tokens, err := h.service.Refresh(r.Context(), req.RefreshToken)
	if err != nil {
		response.Error(w, err)
		return
	}

	response.JSON(w, http.StatusOK, dto.AuthResponse{Token: tokens})
}

// @Summary      Выход из системы
// @Description  Logout пользователя (stateless)
// @Tags         auth
// @Produce      json
// @Success      200   {object} map[string]any
// @Router       /auth/logout [post]
func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	response.JSON(w, http.StatusOK, map[string]any{
		"message": "logged out",
	})
}
