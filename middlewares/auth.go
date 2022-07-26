package middlewares

import (
	"golang_rest_websockets/models"
	"golang_rest_websockets/server"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt"
)

var (
	NO_AUTH_NEEDED = []string{
		"login",
		"signup",
	}
)

func shouldCheckToken(route string) bool {

	for _, noNeededRoute := range NO_AUTH_NEEDED {
		if strings.Contains(route, noNeededRoute) {
			return false
		}
	}
	return true

}

func CheckAuthMiddleware(server server.Server) func(http http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			if !shouldCheckToken(request.URL.Path) {
				next.ServeHTTP(writer, request)
				return
			}
			tokenString := strings.TrimSpace(request.Header.Get("Authorization"))
			_, err := jwt.ParseWithClaims(tokenString, &models.AppClaims{}, func(token *jwt.Token) (interface{}, error) {
				return []byte(server.Config().JWTSecret), nil
			})
			if err != nil {
				http.Error(writer, err.Error(), http.StatusUnauthorized)
				return
			}
			next.ServeHTTP(writer, request)
		})
	}
}
