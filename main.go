package main

import (
	"github.com/Valgard/godotenv"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/justsorrent/game-planner/handlers"
	"log"
	"net/http"
	"os"
)

func main() {
	dotEnv := godotenv.New()
	if err := dotEnv.Load(".env"); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	router := chi.NewRouter()
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Get("/healtz", handlers.HandleHealtz)

	serverPort := os.Getenv("PORT")
	if serverPort == "" {
		log.Fatal("Error loading PORT from .env file")
	}

	log.Printf("Starting server on port %v\n", serverPort)
	log.Fatal(http.ListenAndServe(":"+serverPort, router))
}
