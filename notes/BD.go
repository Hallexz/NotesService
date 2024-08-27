package notes

import (
	"database/sql"
	"log"
	"time"

	_ "github.com/lib/pq"
)

// Глобальная переменная для логирования
var logger *log.Logger

// Функция для установки соединения с базой данных
func SetupDatabase(l *log.Logger) (*sql.DB, error) {
	// Инициализия логгера
	logger = l
	logger.Println("Setting up database connection")

	// Строка подключения к базе данных
	connStr := "user=postgres dbname=notes sslmode=disable"

	// Устанавливаем соединение с базой данных
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		logger.Printf("Failed to open database: %v", err)
		return nil, err
	}

	// Если соединение установлено успешно, выводим лог-сообщение
	logger.Println("Database connection established successfully")
	return db, nil
}

// Структура для представления заметки
type Note struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Интерфейс для работы с заметками
type NoteService interface {
	CreateNote(db *sql.DB, userID int, title, content string) (int, error)
}

// Реализация интерфейса NoteService
type NoteServiceImpl struct{}

// Добавление заметок
func (s *NoteServiceImpl) CreateNote(db *sql.DB, userID int, title, content string) (int, error) {
	// Переменная для хранения ID созданной заметки
	var noteID int

	// Выполняем запрос для создания заметки
	err := db.QueryRow(`
        INSERT INTO notes (user_id, title, content, created_at, updated_at)
        VALUES ($1, $2, $3, NOW(), NOW())
        RETURNING id`, userID, title, content).Scan(&noteID)
	if err != nil {
		// Если ошибка, возвращаем ошибку
		return 0, err
	}
	// Возвращаем ID созданной заметки
	return noteID, nil
}

// Получение списка заметок
func GetNotes(db *sql.DB, userID int) ([]Note, error) {
	// Вывод лог-сообщения о получении заметок
	logger.Printf("Fetching notes for user %d", userID)

	// Выполняем запрос для получения заметок
	rows, err := db.Query(`
        SELECT id, user_id, title, content, created_at, updated_at
        FROM notes
        WHERE user_id = $1`, userID)
	if err != nil {

		// Если ошибка, выводится лог-сообщение и возвращается ошибка
		logger.Printf("Failed to fetch notes: %v", err)
		return nil, err
	}
	defer rows.Close()

	// Переменная для хранения списка заметок
	var notes []Note

	// Проходим по каждой строке результата запроса
	for rows.Next() {

		// Переменная для хранения текущей заметки
		var note Note

		// Считываем данные из строки в структуру Note
		if err := rows.Scan(&note.ID, &note.UserID, &note.Title, &note.Content, &note.CreatedAt, &note.UpdatedAt); err != nil {

			// Если ошибка, выводим лог-сообщение и возвращаем ошибку
			logger.Printf("Failed to scan note: %v", err)
			return nil, err
		}

		// Добавляем заметку в список
		notes = append(notes, note)
	}

	// Выводим лог-сообщение о количестве полученных заметок
	logger.Printf("Retrieved %d notes for user %d", len(notes), userID)

	// Возвращаем список заметок
	return notes, nil
}
