package main

import (
	"github.com/Valgard/godotenv"
	"log"
)

func main() {
	dotEnv := godotenv.New()
	if err := dotEnv.Load(".env"); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
}
