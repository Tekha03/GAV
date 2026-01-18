package dbserver

import (
	"gav/auth"
	"gav/user"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	users user.Repository
}

func NewAuthHandler(users user.Repository) *AuthHandler {
	return &AuthHandler{
		users: users,
	}
}

func (ah *AuthHandler) Register(c *gin.Context) {
	var dto RegisterDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hashedPassword, err := auth.HashPassword(dto.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "password hashing failed"})
		return
	}

	u := user.NewUser(
		0,
		nil,
		nil,
		&user.UserSettings{
			Email: dto.Email,
			PasswordHash: hashedPassword,
			CreatedAt: time.Now(),
		},
		nil,
	)

	if err := ah.users.Create(c.Request.Context(), u); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "cannot create user: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "registered",
		"user_id": u.ID,
	})
}

func (ah *AuthHandler) Login(c *gin.Context) {
	var dto LoginDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	u, err := ah.users.GetByEmail(c, dto.Email)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid email"})
		return
	}

	if !auth.CheckPassword(dto.Password, u.Settings.PasswordHash) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid password"})
		return
	}

	token, err := auth.GenerateToken(int(u.ID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "cannot create token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}
