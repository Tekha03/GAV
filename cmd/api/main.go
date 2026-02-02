package main

import (
	"log"
	"net/http"

	"gav/internal/auth"
	"gav/migrations"
	"gav/storage/sqlite"

	"github.com/go-chi/chi/v5"
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

	r := chi.NewRouter()
	r.Post("/register", authHandler.Register)
	r.Post("/login", authHandler.Login)

	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatal(err)
	}
}
