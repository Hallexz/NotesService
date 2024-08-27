package tests

import (
	"NotesService/speller"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// TestCheckSpelling tests the CheckSpelling function
func TestCheckSpelling(t *testing.T) {
	logger := log.New(ioutil.Discard, "", log.LstdFlags)

	inputText := "Привит"
	expectedText := "Привет"

	mockResponse := []speller.SpellCheckResult{
		{
			Code: 1,
			Pos:  0,
			Len:  6,
			Word: "Привит",
			S:    []string{"Привет"},
		},
	}

	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer mockServer.Close()

	originalURL := speller.YandexSpellerURL
	speller.YandexSpellerURL = mockServer.URL
	defer func() {
		speller.YandexSpellerURL = originalURL
	}()

	result, err := speller.CheckSpelling(inputText, logger)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if strings.TrimSpace(result) != expectedText {
		t.Errorf("Expected '%s', got '%s'", expectedText, result)
	}
}
