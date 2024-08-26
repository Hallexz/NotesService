package main

import (
	"NotesService/auntification"
	"NotesService/notes"
	"log"
	"net/http"
	"os"
)

func setupLogger() *log.Logger {
	logFile, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal("Failed to open log file:", err)
	}
	return log.New(logFile, "", log.Ldate|log.Ltime|log.Lshortfile)
}

func main() {
	logger := setupLogger()

	db, err := notes.SetupDatabase(logger)
	if err != nil {
		logger.Fatal("Failed to setup database:", err)
	}
	defer db.Close()

	auntification.SetLogger(logger)
	authService := auntification.NewAuthService(db)

	http.HandleFunc("/auth", authService.AuthenticateHandler)
	http.Handle("/notes", auntification.JWTAuthMiddleware(http.HandlerFunc(notes.CreateNoteHandler(db, logger))))

	logger.Println("Server is running on :9080")
	logger.Fatal(http.ListenAndServe(":9080", nil))
}
