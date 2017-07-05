package utils

import (
	"encoding/json"
	"net/http"
	"os"
	"strconv"
)

const (
	wordGameDictionaryPathPrefix  string = `https://www.wordgamedictionary.com/word-lists/`
	wordGameDictionaryPathSuffix1 string = `-letter-words/`
	wordGameDictionaryPathSuffix2 string = `-letter-words.json`
	cacheFolderName               string = `../data/`
	cacheNameSuffix               string = `-letterwords.json`
)

//WordResponse - Struct to map the exhaustive list of words from wordgamedictionary.com
type WordResponse struct {
	Word string `json:"word"`
}

//GetWordsFromWordGameDict - Get all possible words of particular length from the wordgamedictionary.com repository
func GetWordsFromWordGameDict(wordLen int) map[string]int {
	var words map[string]int
	var allWords []WordResponse
	words = make(map[string]int)
	if wordLen > 1 && wordLen < 13 {
		path := cacheFolderName + strconv.Itoa(wordLen) + cacheNameSuffix
		if _, err := os.Stat(path); err == nil {
			//Cache file containing all words for this word length exists, so read from cache
			allWords = getJSONDataFromCache(path)
		} else {
			//Cache does not exists, so go get it from the internet
			url := wordGameDictionaryPathPrefix + strconv.Itoa(wordLen) + wordGameDictionaryPathSuffix1 + strconv.Itoa(wordLen) + wordGameDictionaryPathSuffix2
			allWords = getJSONDataOnline(url)
		}
		for _, wordResp := range allWords {
			words[wordResp.Word] = 1
		}

	}
	return words
}

func getJSONDataOnline(url string) []WordResponse {
	req, _ := http.NewRequest("GET", url, nil)
	client := &http.Client{}
	resp, err := client.Do(req)
	Check(err)

	defer resp.Body.Close()

	var allWords []WordResponse
	err = json.NewDecoder(resp.Body).Decode(&allWords)
	Check(err)

	return allWords
}

func getJSONDataFromCache(path string) []WordResponse {
	f, _ := os.Open(path)

	var allWords []WordResponse
	err := json.NewDecoder(f).Decode(&allWords)

	Check(err)
	return allWords
}
