package auntification

import (
	"github.com/dgrijalva/jwt-go"
	"net/http"
	"strconv"
	"strings"
)

func JWTAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if logger != nil {
			logger.Println("Checking JWT token")
		}
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header is missing", http.StatusUnauthorized)
			return
		}

		tokenString := strings.TrimSpace(strings.TrimPrefix(authHeader, "Bearer"))
		if tokenString == "" {
			http.Error(w, "Bearer token is missing", http.StatusUnauthorized)
			return
		}

		claims := &Claims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte("secret_key"), nil // Замените "secret_key" на ваш секретный ключ
		})

		if err != nil || !token.Valid {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		r.Header.Set("UserID", strconv.Itoa(claims.UserID))
		next.ServeHTTP(w, r)
	})
}
