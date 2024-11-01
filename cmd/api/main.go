package main

import (
	"log"
	"quiz-app/internal/initialize"

	"github.com/joho/godotenv"
)

func LoadConfig() bool {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
		return false
	}
	return true
}

func main() {
	LoadConfig()
	initialize.InitApp()
}
