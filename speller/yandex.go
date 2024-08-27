package speller

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

// URL сервиса Яндекс.Спеллер
var YandexSpellerURL = "https://speller.yandex.net/services/spellservice.json/checkText"

// Структура для представления результата проверки орфографии
type SpellCheckResult struct {
	Code int      `json:"code"`
	Pos  int      `json:"pos"`
	Row  int      `json:"row"`
	Col  int      `json:"col"`
	Len  int      `json:"len"`
	Word string   `json:"word"`
	S    []string `json:"s"`
}

// Функция для проверки орфографии текста
func CheckSpelling(text string, logger *log.Logger) (string, error) {
	logger.Printf("Starting spell check for text: %s", text)

	// Параметрыы запроса
	params := url.Values{}
	params.Add("text", text)
	params.Add("lang", "ru,en")

	// Отправление запроса к сервису Яндекс.Спеллер
	resp, err := http.PostForm(YandexSpellerURL, params)
	if err != nil {
		logger.Printf("Failed to send request to Yandex.Speller: %v", err)
		return "", fmt.Errorf("failed to send request to Yandex.Speller: %w", err)
	}
	defer resp.Body.Close()

	// Чтение тело ответа
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.Printf("Failed to read response body: %v", err)
		return "", fmt.Errorf("failed to read response body: %w", err)
	}

	// Разбор ответ в структуру SpellCheckResult
	var results []SpellCheckResult
	if err := json.Unmarshal(body, &results); err != nil {
		logger.Printf("Failed to unmarshal response: %v", err)
		return "", fmt.Errorf("failed to unmarshal response: %w", err)
	}

	logger.Printf("Received %d spell check results", len(results))

	correctedText := []rune(text)
	textLength := len(correctedText)

	// Проход по каждому результату проверки орфографии
	for i := len(results) - 1; i >= 0; i-- {
		result := results[i]
		if len(result.S) > 0 {
			// Определение позиции ошибки в тексте
			startPos := result.Pos
			endPos := result.Pos + result.Len

			// Проверяем допустимость позиции ошибки
			if startPos < 0 || endPos > textLength || startPos > endPos {
				logger.Printf("Invalid position range for word '%s': [%d:%d], skipping", result.Word, startPos, endPos)
				continue
			}

			// Подготавка исправленного слова
			correctedWord := []rune(result.S[0])
			logger.Printf("Correcting word '%s' to '%s' at position %d", result.Word, string(correctedWord), startPos)

			// Замена слова в тексте
			correctedText = append(correctedText[:startPos], append(correctedWord, correctedText[endPos:]...)...)
			textLength = len(correctedText) // Update text length after modification
		}
	}

	// Преобразоваание исправленного текста в строку
	correctedString := string(correctedText)
	logger.Printf("Spell check completed. Corrected text: %s", correctedString)

	return correctedString, nil
}
