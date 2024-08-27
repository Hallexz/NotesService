package notes

import (
	"database/sql"
	"log"
	"time"

	_ "github.com/lib/pq"
)

var logger *log.Logger

func SetupDatabase(l *log.Logger) (*sql.DB, error) {
	logger = l
	logger.Println("Setting up database connection")
	connStr := "user=postgres dbname=notes sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		logger.Printf("Failed to open database: %v", err)
		return nil, err
	}
	logger.Println("Database connection established successfully")
	return db, nil
}

type Note struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type NoteService interface {
	CreateNote(db *sql.DB, userID int, title, content string) (int, error)
}

type NoteServiceImpl struct{}

// Добавление заметок
func (s *NoteServiceImpl) CreateNote(db *sql.DB, userID int, title, content string) (int, error) {
	var noteID int
	err := db.QueryRow(`
        INSERT INTO notes (user_id, title, content, created_at, updated_at)
        VALUES ($1, $2, $3, NOW(), NOW())
        RETURNING id`, userID, title, content).Scan(&noteID)
	if err != nil {
		return 0, err
	}
	return noteID, nil
}

// Получение списка заметок
func GetNotes(db *sql.DB, userID int) ([]Note, error) {
	logger.Printf("Fetching notes for user %d", userID)
	rows, err := db.Query(`
        SELECT id, user_id, title, content, created_at, updated_at
        FROM notes
        WHERE user_id = $1`, userID)
	if err != nil {
		logger.Printf("Failed to fetch notes: %v", err)
		return nil, err
	}
	defer rows.Close()

	var notes []Note
	for rows.Next() {
		var note Note
		if err := rows.Scan(&note.ID, &note.UserID, &note.Title, &note.Content, &note.CreatedAt, &note.UpdatedAt); err != nil {
			logger.Printf("Failed to scan note: %v", err)
			return nil, err
		}
		notes = append(notes, note)
	}
	logger.Printf("Retrieved %d notes for user %d", len(notes), userID)
	return notes, nil
}
