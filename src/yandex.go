package src

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
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

func CheckSpelling(text string) (string, error) {
	params := url.Values{}
	params.Add("text", text)
	params.Add("lang", "ru,en") // Проверка на русском и английском языках

	resp, err := http.PostForm(yandexSpellerURL, params)
	if err != nil {
		return "", fmt.Errorf("failed to send request to Yandex.Speller: %w", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %w", err)
	}

	var results []SpellCheckResult
	if err := json.Unmarshal(body, &results); err != nil {
		return "", fmt.Errorf("failed to unmarshal response: %w", err)
	}

	correctedText := text
	for i := len(results) - 1; i >= 0; i-- {
		result := results[i]
		if len(result.S) > 0 {
			correctedText = correctedText[:result.Pos] + result.S[0] + correctedText[result.Pos+result.Len:]
		}
	}

	return correctedText, nil
}
