package notes

import (
	"NotesService/speller"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
)

type CreateNoteRequest struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

func CreateNoteHandler(db *sql.DB, logger *log.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger.Printf("Received request to create note")

		if r.Method != http.MethodPost {
			logger.Printf("Method not allowed: %s", r.Method)
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		userID := r.Header.Get("UserID")
		if userID == "" {
			logger.Printf("Unauthorized request: missing UserID")
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		var req CreateNoteRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			logger.Printf("Failed to decode request body: %v", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		userIDInt, err := strconv.Atoi(userID)
		if err != nil {
			logger.Printf("Invalid user ID: %s", userID)
			http.Error(w, "Invalid user ID", http.StatusBadRequest)
			return
		}

		logger.Printf("Checking spelling for title")
		correctedTitle, err := speller.CheckSpelling(req.Title, logger)
		if err != nil {
			logger.Printf("Failed to check spelling in title: %v", err)
			http.Error(w, "Failed to check spelling in title: "+err.Error(), http.StatusInternalServerError)
			return
		}
		req.Title = correctedTitle

		logger.Printf("Checking spelling for content")
		correctedContent, err := speller.CheckSpelling(req.Content, logger)
		if err != nil {
			logger.Printf("Failed to check spelling in content: %v", err)
			http.Error(w, "Failed to check spelling in content: "+err.Error(), http.StatusInternalServerError)
			return
		}
		req.Content = correctedContent

		logger.Printf("Creating note in database")
		noteID, err := CreateNote(db, userIDInt, req.Title, req.Content)
		if err != nil {
			logger.Printf("Failed to create note: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		logger.Printf("Note created successfully with ID: %d", noteID)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]int{"id": noteID})
	}
}
