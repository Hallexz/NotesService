package tests

import (
	"NotesService/auntification"
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
)

func TestAuthService_Authenticate(t *testing.T) {
	// Создаем mock базы данных
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	// Создаем экземпляр AuthService с mock базой данных
	authService := auntification.NewAuthService(db)

	// Определяем ожидаемые результаты
	expectedUser := &auntification.User{
		ID:       1,
		Username: "testuser",
		Password: "testpassword",
	}

	// Определяем ожидаемый запрос к базе данных
	rows := sqlmock.NewRows([]string{"id", "username", "password"}).
		AddRow(expectedUser.ID, expectedUser.Username, expectedUser.Password)
	mock.ExpectQuery("SELECT id, username, password FROM users WHERE username = ?").
		WithArgs("testuser").
		WillReturnRows(rows)

	// Вызываем функцию Authenticate
	user, err := authService.Authenticate("testuser", "testpassword")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	// Проверяем результат
	if user.ID != expectedUser.ID || user.Username != expectedUser.Username || user.Password != expectedUser.Password {
		t.Errorf("Unexpected user: %v", user)
	}

	// Проверяем, что все ожидаемые вызовы были выполнены
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There were unfulfilled expectations: %s", err)
	}
}

func TestAuthService_AuthenticateInvalidCredentials(t *testing.T) {
	// Создаем mock базы данных
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	// Создаем экземпляр AuthService с mock базой данных
	authService := auntification.NewAuthService(db)

	// Определяем ожидаемый запрос к базе данных
	rows := sqlmock.NewRows([]string{"id", "username", "password"}).
		AddRow(1, "testuser", "wrongpassword")
	mock.ExpectQuery("SELECT id, username, password FROM users WHERE username = ?").
		WithArgs("testuser").
		WillReturnRows(rows)

	// Вызываем функцию Authenticate с неверным паролем
	_, err = authService.Authenticate("testuser", "testpassword")
	if err == nil {
		t.Error("Expected error for invalid credentials, but got nil")
	} else if err.Error() != "неизвестный паароль" {
		t.Errorf("Unexpected error message: %s", err.Error())
	}

	// Проверяем, что все ожидаемые вызовы были выполнены
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There were unfulfilled expectations: %s", err)
	}
}

func TestAuthService_AuthenticateUserNotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	// Создаем экземпляр AuthService с mock базой данных
	authService := auntification.NewAuthService(db)

	// Определяем ожидаемый запрос к базе данных
	rows := sqlmock.NewRows([]string{"id", "username", "password"})
	mock.ExpectQuery("SELECT id, username, password FROM users WHERE username = ?").
		WithArgs("nonexistentuser").
		WillReturnRows(rows)

	// Вызываем функцию Authenticate с несуществующим пользователем
	_, err = authService.Authenticate("nonexistentuser", "testpassword")
	if err == nil {
		t.Error("Expected error for non-existent user, but got nil")
	} else if err.Error() != "неизвестное имя или паароль" {
		t.Errorf("Unexpected error message: %s", err.Error())
	}

	// Проверяем, что все ожидаемые вызовы были выполнены
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There were unfulfilled expectations: %s", err)
	}
}
