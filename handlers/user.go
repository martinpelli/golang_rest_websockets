package handlers

import (
	"encoding/json"
	"golang_rest_websockets/models"
	"golang_rest_websockets/repositorys"
	"golang_rest_websockets/server"
	"net/http"

	"github.com/segmentio/ksuid"
)

type SingUpRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type SingUpResponse struct {
	Id    string `json:"id"`
	Email string `json:"email"`
}

func SingUPHandler(server server.Server) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		var signUpRequest = SingUpRequest{}
		err := json.NewDecoder(request.Body).Decode(&signUpRequest)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusBadRequest)
			return
		}
		id, err := ksuid.NewRandom()
		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
		}
		var user = models.User{
			Email:    signUpRequest.Email,
			Password: signUpRequest.Password,
			Id:       id.String(),
		}
		err = repositorys.InsertUser(request.Context(), &user)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}
		writer.Header().Set("Content-Type", "application/json")
		json.NewEncoder(writer).Encode(SingUpResponse{
			Id:    user.Id,
			Email: user.Email,
		})
	}
}
