package auth

import (
	"gav/dbserver"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService Service
}

func NewAuthHandler(authService Service) *AuthHandler {
	return &AuthHandler{authService: authService}
}

func (ah *AuthHandler) Register(c *gin.Context) {
	var dto dbserver.RegisterDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	registeredID, err := ah.authService.Register(
		c.Request.Context(),
		dto.Email,
		dto.Password,
	)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user_id": registeredID,
	})
}

func (ah *AuthHandler) Login(c *gin.Context) {
	var dto dbserver.LoginDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	authorizedToken, err := ah.authService.Login(
		c.Request.Context(),
		dto.Email,
		dto.Password,
	)

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": authorizedToken})
}
