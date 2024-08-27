package main

import (
	"NotesService/auntification"
	"NotesService/notes"
	"log"
	"net/http"
	"os"
)

// Функция для настройки логера
func setupLogger() *log.Logger {
	logFile, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal("Failed to open log file:", err)
	}
	return log.New(logFile, "", log.Ldate|log.Ltime|log.Lshortfile)
}

func main() {
	// Настройка логера
	logger := setupLogger()

	// Настройка базы данных
	db, err := notes.SetupDatabase(logger)
	if err != nil {
		logger.Fatal("Failed to setup database:", err)
	}
	defer db.Close()

	// Настройка сервиса аутентификации
	auntification.SetLogger(logger)
	authService := auntification.NewAuthService(db)

	// Создание экземпляра NoteServiceImpl
	noteService := &notes.NoteServiceImpl{}

	// Настройка маршрутов
	http.HandleFunc("/auth", authService.AuthenticateHandler)
	http.Handle("/notes", auntification.JWTAuthMiddleware(notes.CreateNoteHandler(db, logger, noteService)))
	http.Handle("/getNotes", auntification.JWTAuthMiddleware(notes.GetNotesHandler(db)))

	logger.Println("Server is running on :9080")
	logger.Fatal(http.ListenAndServe(":9080", nil))
}
