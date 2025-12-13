package main

import (
	"gav/dbserver"

	"github.com/gin-gonic/gin"
)

func main() {
	dbserver.InitDB()

	r := gin.Default()

	r.POST("/register", dbserver.RegisterHandler)
	r.POST("/login", dbserver.LoginHandler)

	r.Run(":8080")
}
