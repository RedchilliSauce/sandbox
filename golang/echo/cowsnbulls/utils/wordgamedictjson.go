package utils

import (
	"encoding/json"
	"net/http"
	"strconv"
)

const (
	wordGameDictionaryPathPrefix  string = `https://www.wordgamedictionary.com/word-lists/`
	wordGameDictionaryPathSuffix1 string = `-letter-words/`
	wordGameDictionaryPathSuffix2 string = `-letter-words.json`
)

//WordResponse ...
type WordResponse struct {
	Word string `json:"word"`
}

//GetWordsFromWordGameDict ...
func GetWordsFromWordGameDict(wordLen int) map[string]int {
	var words map[string]int
	words = make(map[string]int)
	if wordLen > 1 && wordLen < 13 {
		url := wordGameDictionaryPathPrefix + strconv.Itoa(wordLen) + wordGameDictionaryPathSuffix1 + strconv.Itoa(wordLen) + wordGameDictionaryPathSuffix2
		req, _ := http.NewRequest("GET", url, nil)
		client := &http.Client{}
		resp, err := client.Do(req)

		Check(err)

		defer resp.Body.Close()

		var allWords []WordResponse
		err = json.NewDecoder(resp.Body).Decode(&allWords)
		Check(err)
		for _, wordResp := range allWords {
			words[wordResp.Word] = 1
		}

	}
	return words
}
