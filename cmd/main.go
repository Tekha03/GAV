package main

import (
	"gav/dbserver"

	"github.com/gin-gonic/gin"
)

func main() {
	dbserver.InitDB()

	r := gin.Default()

	r.POST("/register", dbserver.RegisterHandler)
	r.GET("/users/:id", dbserver.GetUserHandler)

	r.Run(":8080")
}
