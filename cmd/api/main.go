package main

import (
	"log"

	"gav/internal/auth"
	"gav/migrations"
	"gav/storage/sqlite"

	"github.com/gin-gonic/gin"
)

func main() {
	db, err := sqlite.Open("social.db")
	if err != nil {
		log.Fatal(err)
	}

	if err := migrations.Migrate(db); err != nil {
		log.Fatal(err)
	}

	userRepo := sqlite.NewUserRepository(db)
	// postRepo := memory.NewPostRepository()
	authService := auth.NewAuthService(userRepo)
	authHandler := auth.NewAuthHandler(authService)

	r := gin.Default()
	r.POST("/register", authHandler.Register)
	r.POST("/login", authHandler.Login)

	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
