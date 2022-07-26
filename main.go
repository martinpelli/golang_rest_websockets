package main

import (
	"context"
	"golang_rest_websockets/handlers"
	"golang_rest_websockets/middlewares"
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
	router.Use(middlewares.CheckAuthMiddleware(server))

	router.HandleFunc("/", handlers.HomeHandler(server)).Methods(http.MethodGet)
	router.HandleFunc("/signup", handlers.SingUpLoginHandler(server)).Methods(http.MethodPost)
	router.HandleFunc("/login", handlers.LoginHandler(server)).Methods(http.MethodPost)
	router.HandleFunc("/me", handlers.MeHandler(server)).Methods(http.MethodGet)
	router.HandleFunc("/posts", handlers.InsertPostHandler(server)).Methods(http.MethodPost)
	router.HandleFunc("/posts/{id}", handlers.GetPostByIdHandler(server)).Methods(http.MethodGet)
	router.HandleFunc("/posts/{id}", handlers.UpdatePostHandler(server)).Methods(http.MethodPut)
	router.HandleFunc("/posts/{id}", handlers.DeletePostHandler(server)).Methods(http.MethodDelete)
	router.HandleFunc("/posts", handlers.ListPostHandler(server)).Methods(http.MethodGet)
	router.HandleFunc("/ws", server.Hub().HandleWebSocket)
}
