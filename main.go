package main

import (
	"database/sql"
	"github.com/Valgard/godotenv"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/justsorrent/game-planner/handlers"
	"github.com/justsorrent/game-planner/internal/db"
	"log"
	"net/http"
	"os"

	_ "github.com/lib/pq"
)

func main() {
	dotEnv := godotenv.New()
	if err := dotEnv.Load(".env"); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	dbString := os.Getenv("DB_STRING")
	if dbString == "" {
		log.Fatal("Error loading DB_STRING from .env file")
	}

	conn, err := sql.Open("postgres", dbString)
	if err != nil {
		log.Fatalf("Cannot connect to DB. %v", err)
	}

	cfg := handlers.ApiConfig{
		DB: db.New(conn),
	}

	router := chi.NewRouter()
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Get("/healtz", handlers.HandleHealtz)

	v1Router := chi.NewRouter()
	// games routes
	// TODO: add auth middleware
	v1Router.Get("/games", cfg.HandleGetGames)
	v1Router.Post("/games", cfg.HandleCreateGame)
	v1Router.Put("/games/{id}", cfg.HandleUpdateGameById)
	v1Router.Get("/games/{id}", cfg.HandleGetGameById)
	v1Router.Delete("/games/{id}", cfg.HandleDeleteGameById)
	// users routes
	v1Router.Post("/users", cfg.HandleCreateUser)
	// session routes
	v1Router.Post("/session", cfg.HandleCreateSession)
	v1Router.Put("/session", cfg.HandleRefreshSession)
	v1Router.Delete("/session", cfg.HandleDeleteSession)
	router.Mount("/v1", v1Router)

	serverPort := os.Getenv("PORT")
	if serverPort == "" {
		log.Fatal("Error loading PORT from .env file")
	}

	log.Printf("Starting server on port %v\n", serverPort)
	log.Fatal(http.ListenAndServe(":"+serverPort, router))
}
