package handlers

import (
	"encoding/json"
	"golang_rest_websockets/models"
	"golang_rest_websockets/repositorys"
	"golang_rest_websockets/server"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/segmentio/ksuid"
	"golang.org/x/crypto/bcrypt"
)

const (
	HASH_COST = 8
)

type SingUpRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type SingUpResponse struct {
	Id    string `json:"id"`
	Email string `json:"email"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

func SingUpLoginHandler(server server.Server) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		var signUpRequest = SingUpRequest{}
		err := json.NewDecoder(request.Body).Decode(&signUpRequest)

		if err != nil {
			http.Error(writer, err.Error(), http.StatusBadRequest)
			return
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(signUpRequest.Password), HASH_COST)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}

		id, err := ksuid.NewRandom()
		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
		}
		var user = models.User{
			Email:    signUpRequest.Email,
			Password: string(hashedPassword),
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

func LoginHandler(server server.Server) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		var loginRequest = SingUpRequest{}
		err := json.NewDecoder(request.Body).Decode(&request)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusBadRequest)
			return
		}

		user, err := repositorys.GetUserByEmail(request.Context(), loginRequest.Email)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}
		if user == nil {
			http.Error(writer, "Invalid credentials", http.StatusUnauthorized)
		}

		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(user.Password)); err != nil {
			http.Error(writer, "Invalid credentials", http.StatusUnauthorized)
			return
		}

		claims := models.AppClaims{
			UserId: user.Id,
			StandardClaims: jwt.StandardClaims{
				ExpiresAt: time.Now().Add(2 * time.Hour * 24).Unix(),
			},
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		signedToken, err := token.SignedString([]byte(server.Config().JWTSecret))
		if user != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}

		writer.Header().Set("Content-Type", "application/json")
		json.NewEncoder(writer).Encode(LoginResponse{
			Token: signedToken,
		})

	}
}
