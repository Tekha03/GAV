package main

import (
	"log"

	"gav/internal/auth"
	"gav/internal/post"
	"gav/storage/sqlite"
	"gav/user"

	"github.com/gin-gonic/gin"
)

func main() {
	db, err := sqlite.Open("social.db")
	if err != nil {
		log.Fatal(err)
	}

	if err := db.AutoMigrate(
		&user.User{},
		&post.Post{},
	); err != nil {
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
