package main

import (
	"context"
	"log"

	"gav/internal/app"
	"gav/internal/config"
)

func main() {
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
