package tests

import (
	"NotesService/notes"
	"database/sql"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// Mock for Speller
type MockNoteService struct {
	mock.Mock
}

func (m *MockNoteService) CreateNote(db *sql.DB, userID int, title, content string) (int, error) {
	args := m.Called(db, userID, title, content)
	return args.Int(0), args.Error(1)
}

func TestCreateNoteHandler(t *testing.T) {
	// Create mock logger
	logger := log.New(log.Writer(), "test: ", log.LstdFlags)

	// Create a mock database (we won't actually use it in this test)
	db := &sql.DB{}

	// Create a mock NoteService
	mockNoteService := new(MockNoteService)
	mockNoteService.On("CreateNote", mock.Anything, 1, "Test Notebooks", "This is a test note").Return(1, nil)

	// Create the handler
	handler := notes.CreateNoteHandler(db, logger, mockNoteService)

	// Create a new HTTP request
	reqBody := `{"title":"Test Note","content":"This is a test note"}`
	req := httptest.NewRequest(http.MethodPost, "/notes", strings.NewReader(reqBody))
	req.Header.Set("UserID", "1")

	// Create a ResponseRecorder to capture the response
	rr := httptest.NewRecorder()

	// Call the handler
	handler.ServeHTTP(rr, req)

	// Check the response status code
	assert.Equal(t, http.StatusCreated, rr.Code)

	// Parse the response body
	var respBody map[string]int
	err := json.NewDecoder(rr.Body).Decode(&respBody)
	assert.NoError(t, err)

	// Check that the note ID is as expected
	assert.Equal(t, 1, respBody["id"])

	// Assert that the mock service was called
	mockNoteService.AssertExpectations(t)
}
