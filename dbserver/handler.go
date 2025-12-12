package dbserver

import (
	"gav/user"
	"net/http"

	"github.com/gin-gonic/gin"
)

func RegisterHandler(c *gin.Context) {
	var u user.User
	if err := c.ShouldBindJSON(&u); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	DB.Create(&u)
	c.JSON(http.StatusOK, u)
}

func GetUserHandler(c *gin.Context) {
	var u user.User
	id := c.Param("id")

	if err := DB.First(&u, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, u)
}
