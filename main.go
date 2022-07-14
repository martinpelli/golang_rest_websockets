package main

import (
	"context"
	"golang_rest_websockets/handlers"
	"golang_rest_websockets/server"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading environment: .env file not loaded")
	}

	PORT := os.Getenv("PORT")
	JWT_SECRET := os.Getenv("JWT_SECRET")
	DATABASE_URL := os.Getenv("DATABASE_URL")

	server, err := server.NewServer(context.Background(), &server.Config{
		JWTSecret:   JWT_SECRET,
		Port:        PORT,
		DatabaseUrl: DATABASE_URL,
	})

	if err != nil {
		log.Fatalf("Error creating server: %v", err)
	}

	server.Start(BindRoutes)
}

func BindRoutes(server server.Server, router *mux.Router) {
	router.HandleFunc("/", handlers.HomeHandler(server)).Methods(http.MethodGet)
	router.HandleFunc("/signup", handlers.SingUPHandler(server)).Methods(http.MethodPost)
}
