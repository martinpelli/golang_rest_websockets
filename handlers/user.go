package handlers

import (
	"encoding/json"
	"golang_rest_websockets/models"
	"golang_rest_websockets/repositorys"
	"golang_rest_websockets/server"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/segmentio/ksuid"
	"golang.org/x/crypto/bcrypt"
)

const (
	HASH_COST = 8
)

type SingUpLoginRequest struct {
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
		var signUpRequest = SingUpLoginRequest{}
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
		var loginRequest = SingUpLoginRequest{}
		err := json.NewDecoder(request.Body).Decode(&loginRequest)
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
		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginRequest.Password)); err != nil {
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
		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}

		writer.Header().Set("Content-Type", "application/json")
		json.NewEncoder(writer).Encode(LoginResponse{
			Token: signedToken,
		})

	}
}

func MeHandler(server server.Server) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		tokenString := strings.TrimSpace(request.Header.Get("Authorization"))
		token, err := jwt.ParseWithClaims(tokenString, &models.AppClaims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(server.Config().JWTSecret), nil
		})
		if err != nil {
			http.Error(writer, err.Error(), http.StatusUnauthorized)
			return
		}

		if claims, ok := token.Claims.(*models.AppClaims); ok && token.Valid {
			user, err := repositorys.GetUserById(request.Context(), claims.UserId)
			if err != nil {
				http.Error(writer, err.Error(), http.StatusInternalServerError)
				return
			}

			writer.Header().Set("Content-Type", "application/json")
			json.NewEncoder(writer).Encode(user)
		} else {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
