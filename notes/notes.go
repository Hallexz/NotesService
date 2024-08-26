package notes

import (
	"NotesService/src"
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
)

type CreateNoteRequest struct {
	UserID  int    `json:"user_id"`
	Title   string `json:"title"`
	Content string `json:"content"`
}

func CreateNoteHandler(db *sql.DB, w http.ResponseWriter, r *http.Request) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		userID := r.Header.Get("UserID")
		if userID == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		var req CreateNoteRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Проверка орфографии заголовка
		correctedTitle, err := src.CheckSpelling(req.Title)
		if err != nil {
			http.Error(w, "Failed to check spelling in title: "+err.Error(), http.StatusInternalServerError)
			return
		}
		req.Title = correctedTitle

		// Проверка орфографии содержимого
		correctedContent, err := src.CheckSpelling(req.Content)
		if err != nil {
			http.Error(w, "Failed to check spelling in content: "+err.Error(), http.StatusInternalServerError)
			return
		}
		req.Content = correctedContent

		// Преобразование userID из строки в int
		userIDInt, err := strconv.Atoi(userID)
		if err != nil {
			http.Error(w, "Invalid user ID", http.StatusBadRequest)
			return
		}

		// Создание заметки с исправленным текстом
		noteID, err := CreateNote(db, userIDInt, req.Title, req.Content)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]int{"id": noteID})
	}
}
