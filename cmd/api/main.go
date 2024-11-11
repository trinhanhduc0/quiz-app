package main

import (
	"fmt"
	"log"
	"os"
	"quiz-app/internal/initialize"

	"github.com/joho/godotenv"
)

func LoadConfig() bool {
	// Không cần dòng này khi triển khai trên Fly.io
	err := godotenv.Load(".env")
	if err != nil {

	port := os.Getenv("PORT")
	appEnv := os.Getenv("APP_ENV")
	redisURI := os.Getenv("REDIS_URI")
	mongoDBURI := os.Getenv("MONGODB_URI")
	awsAccessKeyID := os.Getenv("aws_access_key_id")
	awsSecretAccessKey := os.Getenv("aws_secret_access_key")

	// Kiểm tra nếu các biến môi trường quan trọng có giá trị hay không
	if port == "" || appEnv == "" || redisURI == "" || mongoDBURI == "" || awsAccessKeyID == "" || awsSecretAccessKey == "" {
		fmt.Println(port)
		log.Fatal("One or more environment variables are not set.")
		return false
	}
	
	return true
}

func main() {
	if !LoadConfig() {
		//Error load env file
		return
	}
	initialize.InitApp()
}
