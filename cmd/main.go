package main

import (
	"log"

	"gav/dbserver"
	"gav/post"
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
	authHandler := dbserver.NewAuthHandler(userRepo)

	r := gin.Default()
	r.POST("/register", authHandler.Register)
	r.POST("/login", authHandler.Login)

	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
