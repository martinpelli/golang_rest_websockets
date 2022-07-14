package handlers

import (
	"encoding/json"
	"golang_rest_websockets/server"
	"net/http"
)

type HomeReponse struct {
	Message string `json:"message"`
	Status  bool   `json:"status"`
}

func HomeHandler(server server.Server) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusOK)
		json.NewEncoder(writer).Encode(HomeReponse{Message: "Welcome to Golang", Status: true})
	}
}
