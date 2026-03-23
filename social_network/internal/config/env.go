package config

import (
	"log"
	"path/filepath"
	"runtime"

	"github.com/joho/godotenv"
)

func loadEnv() {
	_, filename, _, _:= runtime.Caller(0)
	root := filepath.Join(filepath.Dir(filename), "..", "..")

	envPath := filepath.Join(root, ".env")
	err := godotenv.Load(envPath)
	if err != nil {
		log.Printf("Warning: could not load .env from %s: %v", envPath, err)
	} else {
		log.Printf(".env loaded successfully from %s", envPath)
	}
}
