package speller

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

const yandexSpellerURL = "https://speller.yandex.net/services/spellservice.json/checkText"

type SpellCheckResult struct {
	Code int      `json:"code"`
	Pos  int      `json:"pos"`
	Row  int      `json:"row"`
	Col  int      `json:"col"`
	Len  int      `json:"len"`
	Word string   `json:"word"`
	S    []string `json:"s"`
}

func CheckSpelling(text string, logger *log.Logger) (string, error) {
	logger.Printf("Starting spell check for text: %s", text)

	params := url.Values{}
	params.Add("text", text)
	params.Add("lang", "ru,en") // Проверка на русском и английском языках

	resp, err := http.PostForm(yandexSpellerURL, params)
	if err != nil {
		logger.Printf("Failed to send request to Yandex.Speller: %v", err)
		return "", fmt.Errorf("failed to send request to Yandex.Speller: %w", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.Printf("Failed to read response body: %v", err)
		return "", fmt.Errorf("failed to read response body: %w", err)
	}

	var results []SpellCheckResult
	if err := json.Unmarshal(body, &results); err != nil {
		logger.Printf("Failed to unmarshal response: %v", err)
		return "", fmt.Errorf("failed to unmarshal response: %w", err)
	}

	logger.Printf("Received %d spell check results", len(results))

	correctedText := []rune(text) // Используем []rune, чтобы корректно работать с юникодными символами
	for i := len(results) - 1; i >= 0; i-- {
		result := results[i]
		if len(result.S) > 0 {
			// Исправляем слово, начиная с позиции Pos и до Pos + Len
			correctedWord := result.S[0]
			startPos := result.Pos
			endPos := result.Pos + result.Len

			logger.Printf("Correcting word '%s' to '%s' at position %d", result.Word, correctedWord, startPos)

			// Заменяем слово
			correctedText = append(correctedText[:startPos], append([]rune(correctedWord), correctedText[endPos:]...)...)
		}
	}

	correctedString := string(correctedText)
	logger.Printf("Spell check completed. Corrected text: %s", correctedString)

	return correctedString, nil
}
