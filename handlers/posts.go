package handlers

import (
	"encoding/json"
	"golang_rest_websockets/models"
	"golang_rest_websockets/repositorys"
	"golang_rest_websockets/server"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt"
	"github.com/segmentio/ksuid"
)

type InsertPostRequest struct {
	PostContent string `json:"post_content"`
}

type PostResponse struct {
	Id          string `json:"id"`
	PostContent string `json:"post_content"`
}

func InsertPostHandler(server server.Server) http.HandlerFunc {
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
			var postRequest = InsertPostRequest{}
			if err := json.NewDecoder(request.Body).Decode(&postRequest); err != nil {
				http.Error(writer, err.Error(), http.StatusBadRequest)
				return
			}
			id, err := ksuid.NewRandom()
			if err != nil {
				http.Error(writer, err.Error(), http.StatusInternalServerError)
				return
			}
			post := models.Post{
				Id:          id.String(),
				PostContent: postRequest.PostContent,
				UserId:      claims.UserId,
			}
			err = repositorys.InsertPost(request.Context(), &post)
			if err != nil {
				http.Error(writer, err.Error(), http.StatusInternalServerError)
				return
			}
			writer.Header().Set("Content-Type", "application/json")
			json.NewEncoder(writer).Encode(PostResponse{
				Id:          post.Id,
				PostContent: post.PostContent,
			})
		} else {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
