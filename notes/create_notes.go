package notes

import (
	"NotesService/speller"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
)

// Структура для представления запроса на создание заметки
type CreateNoteRequest struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

// Функция-обработчик запросов на создание заметок
func CreateNoteHandler(db *sql.DB, logger *log.Logger, noteService NoteService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Выводим лог-сообщение о получении запроса
		logger.Printf("Received request to create note")

		// Проверяется метод запроса
		if r.Method != http.MethodPost {
			logger.Printf("Method not allowed: %s", r.Method)
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Получение ID пользователя из заголовка запроса
		userID := r.Header.Get("UserID")
		if userID == "" {
			logger.Printf("Unauthorized request: missing UserID")
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Декодирование тела запроса в структуру CreateNoteRequest
		var req CreateNoteRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			logger.Printf("Failed to decode request body: %v", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Преобразование ID пользователя в целое число
		userIDInt, err := strconv.Atoi(userID)
		if err != nil {
			logger.Printf("Invalid user ID: %s", userID)
			http.Error(w, "Invalid user ID", http.StatusBadRequest)
			return
		}

		// Проверка орфографию в заголовке заметки
		logger.Printf("Checking spelling for title")
		correctedTitle, err := speller.CheckSpelling(req.Title, logger)
		if err != nil {
			logger.Printf("Failed to check spelling in title: %v", err)
			http.Error(w, "Failed to check spelling in title: "+err.Error(), http.StatusInternalServerError)
			return
		}
		req.Title = correctedTitle

		// Проверяка орфографию в содержании заметки
		logger.Printf("Checking spelling for content")
		correctedContent, err := speller.CheckSpelling(req.Content, logger)
		if err != nil {
			logger.Printf("Failed to check spelling in content: %v", err)
			http.Error(w, "Failed to check spelling in content: "+err.Error(), http.StatusInternalServerError)
			return
		}
		req.Content = correctedContent

		// Создание заметки в базе данных
		logger.Printf("Creating note in database with Title: %s, Content: %s", req.Title, req.Content)
		noteID, err := noteService.CreateNote(db, userIDInt, req.Title, req.Content)
		if err != nil {
			// Если создание заметки прошло неудачно, выводим лог-сообщение и возвращаем ошибку
			logger.Printf("Failed to create note: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Вывод лог-сообщения о успешном создании заметки
		logger.Printf("Note created successfully with ID: %d", noteID)
		// Возвращение ответа клиенту
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]int{"id": noteID})
	}
}
