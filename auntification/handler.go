package auntification

import (
	"encoding/json"
	"github.com/dgrijalva/jwt-go"
	"net/http"
	"time"
)

// Структура для представления данных, содержащихся в JWT-токене
type Claims struct {
	UserID   int    `json:"user_id"`
	Username string `json:"username"`
	jwt.StandardClaims
}

// Функция-обработчик запросов на аутентификацию
func (s *AuthService) AuthenticateHandler(w http.ResponseWriter, r *http.Request) {
	// Выводим лог-сообщение о получении запроса на аутентификацию
	if logger != nil {
		logger.Println("Received authentication request")
	}

	// Структура для представления данных аутентификации
	var credentials struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	// Декодируем тело запроса в структуру credentials
	err := json.NewDecoder(r.Body).Decode(&credentials)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Аутентифицируем пользователя
	user, err := s.Authenticate(credentials.Username, credentials.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	// Создание JWT-токена
	claims := Claims{
		UserID:   user.ID,
		Username: user.Username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// Подписывание токена секретным ключом
	tokenString, err := token.SignedString([]byte("secret_key")) // Замените "secret_key" на свой секретный ключ
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Устанавливаем заголовок ответа для указания формата данных
	w.Header().Set("Content-Type", "application/json")
	// Возвращаем JWT-токен в ответе
	json.NewEncoder(w).Encode(map[string]string{"token": tokenString})
}
