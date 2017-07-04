package utils

import (
	"io/ioutil"
	"strings"
)

//Check ...
func Check(e error) {
	if e != nil {
		panic(e)
	}
}

func GetWordsFromFile(path string, wordLen int) map[string]int {
	var words map[string]int
	if wordLen > 1 && wordLen < 1000 {
		wordsRawData := readFile(path)
		words = getNLetterWords(wordsRawData, wordLen)
	} else {
		words = nil
	}
	return words
}

func readFile(path string) string {
	dat, err := ioutil.ReadFile(path)
	Check(err)
	return string(dat)
}

func getNLetterWords(wordsRawData string, wordLen int) map[string]int {
	words := make(map[string]int)
	allWords := strings.Split(wordsRawData, "\n")
	for _, word := range allWords {
		if len(word) == wordLen {
			words[word] = 1
		}
	}
	return words
}
