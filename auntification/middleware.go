package auntification

import (
	"github.com/dgrijalva/jwt-go"
	"net/http"
	"strconv"
	"strings"
)

// Функция для создания middleware-а для проверки JWT-токена
func JWTAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Выводим лог-сообщение о проверке JWT-токена
		if logger != nil {
			logger.Println("Checking JWT token")
		}

		// Получаем значение заголовка Authorization из запроса
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header is missing", http.StatusUnauthorized)
			return
		}

		// Извлекаем JWT-токен из заголовка Authorization
		tokenString := strings.TrimSpace(strings.TrimPrefix(authHeader, "Bearer"))
		if tokenString == "" {
			http.Error(w, "Bearer token is missing", http.StatusUnauthorized)
			return
		}

		claims := &Claims{}
		// Парсим JWT-токен
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte("secret_key"), nil // Замените "secret_key" на ваш секретный ключ
		})

		if err != nil || !token.Valid {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		// Устанавливаем значение заголовка UserID в запросе
		r.Header.Set("UserID", strconv.Itoa(claims.UserID))
		next.ServeHTTP(w, r)
	})
}
