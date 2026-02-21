package main

import (
	"context"
	"log"

	"gav/dbserver"
	"gav/internal/app"
	"gav/internal/config"
)

func main() {
	db, err := dbserver.InitDB()
	if err != nil {
		log.Fatal(err)
	}

	defer func () {
		if err := dbserver.CloseDB(db); err != nil {
			log.Printf("failed to close db: %v", err)
		}
	}()

	ctx := context.Background()

	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	app, err := app.NewApp(ctx, cfg)
	if err != nil {
		log.Fatal(err)
	}

	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
