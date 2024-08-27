package notes

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
)

func GetNotesHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Получаем userID из параметров запроса
		userIDStr := r.URL.Query().Get("user_id")
		if userIDStr == "" {
			http.Error(w, "Missing user_id parameter", http.StatusBadRequest)
			return
		}

		userID, err := strconv.Atoi(userIDStr)
		if err != nil {
			http.Error(w, "Invalid user_id", http.StatusBadRequest)
			return
		}

		// Получаем заметки из базы данных
		notes, err := GetNotes(db, userID)
		if err != nil {
			http.Error(w, "Failed to fetch notes", http.StatusInternalServerError)
			return
		}

		// Отправляем заметки в формате JSON
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(notes)
	}
}
