package config

import "github.com/joho/godotenv"

func loadEnv() {
	_ = godotenv.Load()
}
