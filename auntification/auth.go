package auntification

import (
	"database/sql"
	"errors"
	"log"
)

var logger *log.Logger

func SetLogger(l *log.Logger) {
	logger = l
}

type AuthService struct {
	db *sql.DB
}

func NewAuthService(db *sql.DB) *AuthService {
	return &AuthService{db: db}
}

func (s *AuthService) Authenticate(username, password string) (*User, error) {
	query := "SELECT id, username, password FROM users WHERE username = $1"
	row := s.db.QueryRow(query, username)

	var user User
	err := row.Scan(&user.ID, &user.Username, &user.Password)
	if err != nil {
		return nil, errors.New("неизвестное имя или паароль")
	}

	if user.Password != password {
		return nil, errors.New("неизвестный паароль")
	}

	return &user, nil
}
