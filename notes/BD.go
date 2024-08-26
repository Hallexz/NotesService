package notes

import (
	"database/sql"
	"log"
	"time"

	_ "github.com/lib/pq"
)

func SetupDatabase() (*sql.DB, error) {
	connStr := "user=postgres dbname=notes sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
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

// Добавление заметок
func CreateNote(db *sql.DB, userID int, title, content string) (int, error) {
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
	rows, err := db.Query(`
		SELECT id, user_id, title, content, created_at, updated_at
		FROM notes
		WHERE user_id = $1`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var notes []Note
	for rows.Next() {
		var note Note
		if err := rows.Scan(&note.ID, &note.UserID, &note.Title, &note.Content, &note.CreatedAt, &note.UpdatedAt); err != nil {
			return nil, err
		}
		notes = append(notes, note)
	}
	return notes, nil
}
