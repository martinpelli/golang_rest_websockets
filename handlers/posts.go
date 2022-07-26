package handlers

import (
	"encoding/json"
	"golang_rest_websockets/models"
	"golang_rest_websockets/repositorys"
	"golang_rest_websockets/server"
	"net/http"
	"strconv"
	"strings"

	"github.com/golang-jwt/jwt"
	"github.com/gorilla/mux"
	"github.com/segmentio/ksuid"
)

type UpsertPostRequest struct {
	PostContent string `json:"postContent"`
}

type PostResponse struct {
	Id          string `json:"id"`
	PostContent string `json:"postContent"`
}

type PostUpdateResponse struct {
	Message string `json:"message"`
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
			var postRequest = UpsertPostRequest{}
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
			var postMessage = models.Websocketmessage{
				Type:    "Post_Created",
				Payload: post,
			}
			server.Hub().Broadcast(postMessage, nil)
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

func GetPostByIdHandler(server server.Server) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		params := mux.Vars(request)
		post, err := repositorys.GestPostById(request.Context(), params["id"])
		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}
		writer.Header().Set("Content-Type", "application/json")
		json.NewEncoder(writer).Encode(post)
	}
}

func UpdatePostHandler(server server.Server) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		params := mux.Vars(request)
		tokenString := strings.TrimSpace(request.Header.Get("Authorization"))
		token, err := jwt.ParseWithClaims(tokenString, &models.AppClaims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(server.Config().JWTSecret), nil
		})
		if err != nil {
			http.Error(writer, err.Error(), http.StatusUnauthorized)
			return
		}

		if claims, ok := token.Claims.(*models.AppClaims); ok && token.Valid {
			var postRequest = UpsertPostRequest{}
			if err := json.NewDecoder(request.Body).Decode(&postRequest); err != nil {
				http.Error(writer, err.Error(), http.StatusBadRequest)
				return
			}
			post := models.Post{
				Id:          params["id"],
				PostContent: postRequest.PostContent,
				UserId:      claims.UserId,
			}
			err = repositorys.UpdatePost(request.Context(), &post)
			if err != nil {
				http.Error(writer, err.Error(), http.StatusInternalServerError)
				return
			}
			writer.Header().Set("Content-Type", "application/json")
			json.NewEncoder(writer).Encode(PostUpdateResponse{
				Message: "Post actualizado",
			})
		} else {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func DeletePostHandler(server server.Server) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		params := mux.Vars(request)
		tokenString := strings.TrimSpace(request.Header.Get("Authorization"))
		token, err := jwt.ParseWithClaims(tokenString, &models.AppClaims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(server.Config().JWTSecret), nil
		})
		if err != nil {
			http.Error(writer, err.Error(), http.StatusUnauthorized)
			return
		}

		if claims, ok := token.Claims.(*models.AppClaims); ok && token.Valid {

			err = repositorys.DeletePost(request.Context(), params["id"], claims.UserId)
			if err != nil {
				http.Error(writer, err.Error(), http.StatusInternalServerError)
				return
			}
			writer.Header().Set("Content-Type", "application/json")
			json.NewEncoder(writer).Encode(PostUpdateResponse{
				Message: "Post borrado",
			})
		} else {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func ListPostHandler(server server.Server) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		var err error
		pageStr := request.URL.Query().Get("page")
		var page = uint64(0)
		if pageStr != "" {
			page, err = strconv.ParseUint(pageStr, 10, 64)
			if err != nil {
				http.Error(writer, err.Error(), http.StatusBadRequest)
				return
			}
		}
		posts, err := repositorys.ListPost(request.Context(), page)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}

		writer.Header().Set("Content-Type", "application/json")
		json.NewEncoder(writer).Encode(posts)
	}
}
