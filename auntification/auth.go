package auntification

import (
	"database/sql"
	"errors"
	"log"
)

// Глобальная переменная для логирования
var logger *log.Logger

// Функция для установки логера
func SetLogger(l *log.Logger) {
	logger = l
}

// Структура для представления сервиса аутентификации
type AuthService struct {
	db *sql.DB
}

// Функция для создания нового сервиса аутентификации
func NewAuthService(db *sql.DB) *AuthService {
	return &AuthService{db: db}
}

// Метод для аутентификации пользователя
func (s *AuthService) Authenticate(username, password string) (*User, error) {
	// Запрос для получения данных пользователя из базы данных
	query := "SELECT id, username, password FROM users WHERE username = $1"

	// Выполняем запрос и получаем результат
	row := s.db.QueryRow(query, username)

	var user User
	// Считывание данных из результата запроса в структуру User
	err := row.Scan(&user.ID, &user.Username, &user.Password)
	if err != nil {
		return nil, errors.New("неизвестное имя или пароль")
	}

	// Проверка пароля
	if user.Password != password {
		return nil, errors.New("неизвестный пароль")
	}

	return &user, nil
}
