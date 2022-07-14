package server

import (
	"errors"
	"golang_rest_websockets/database"
	"golang_rest_websockets/repositorys"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"golang.org/x/net/context"
)

type Config struct {
	Port        string
	JWTSecret   string
	DatabaseUrl string
}

type Server interface {
	Config() *Config
}

type Broker struct {
	config *Config
	router *mux.Router
}

func (broker *Broker) Config() *Config {
	return broker.config
}

func NewServer(context context.Context, config *Config) (*Broker, error) {
	if config.Port == "" {
		return nil, errors.New("Port is empty, port is required")
	}

	if config.JWTSecret == "" {
		return nil, errors.New("Secret Key for JWT is empty, secret key is required")
	}

	if config.DatabaseUrl == "" {
		return nil, errors.New("Database URL is empty, database url is required")
	}

	broker := &Broker{
		config: config,
		router: mux.NewRouter(),
	}

	return broker, nil
}

func (broker *Broker) Start(binder func(server Server, router *mux.Router)) {
	broker.router = mux.NewRouter()
	binder(broker, broker.router)
	repo, err := database.NewPostgresRepository(broker.config.DatabaseUrl)
	if err != nil {
		log.Fatal(err)
	}
	repositorys.SetRepository(repo)
	log.Println("Starting server on port", broker.Config().Port)
	if err := http.ListenAndServe(broker.config.Port, broker.router); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
