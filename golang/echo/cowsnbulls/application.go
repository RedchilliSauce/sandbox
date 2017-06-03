package main

import (
	"log"
	"net/http"
	"os"

	"github.com/labstack/echo"
)

//Game ...
type Game struct {
	Word    string              `json:"word"`
	Guesses map[string]GuessRes `json:"guesses"`
}

//GuessRes ...
type GuessRes struct {
	Cows  int `json:"cows"`
	Bulls int `json:"bulls"`
}

//Games ...
var Games map[string]Game

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	f, _ := os.Create("/var/log/forecast.log")
	defer f.Close()
	log.SetOutput(f)

	const indexPage = "public/index.html"

	Games = make(map[string]Game)
	e := echo.New()

	e.GET("/sendguess/:p1/:p2/:guess", guess)
	e.GET("/creategame/:p1/:p2/:word", creategame)

	e.Start(":" + port)
}

func guess(c echo.Context) error {
	p1 := c.Param("p1")
	p2 := c.Param("p2")
	guess := c.Param("guess")

	key := generateKey(p1, p2)

	game, exists := Games[key]

	if exists {
		res := getGuessRes(game.Word, guess)
		if res == nil {
			return c.String(http.StatusBadRequest, "Check the number of letters in your input")
		}
		if res.Bulls == 4 && res.Cows == 0 {
			return c.String(http.StatusOK, "Whee you guessed the word correctly. Word was *"+game.Word+"*")
		}
		game.Guesses[guess] = *res
	}
	return c.JSON(http.StatusOK, game.Guesses)
}

func generateKey(p1 string, p2 string) string {
	return p1 + "#@#" + p2
}

func creategame(c echo.Context) error {
	guesses := make(map[string]GuessRes)

	p1 := c.Param("p1")
	p2 := c.Param("p2")
	word := c.Param("word")

	key := generateKey(p1, p2)

	newGame := Game{word, guesses}

	_, exists := Games[key]
	if exists {
		return c.String(http.StatusBadRequest, "Game already exists between the two players")
	}

	Games[key] = newGame
	return c.String(http.StatusOK, "Game created")
}

func getGuessRes(expected string, actual string) *GuessRes {
	if len(expected) != len(actual) {
		return nil
	}

	cows := 0
	bulls := 0

	bullTracker := make(map[int]int)
	for i := 0; i < len(expected); i++ {
		if expected[i] == actual[i] {
			bulls++
			bullTracker[i] = 1
		}
	}

	if bulls < 4 {
		for i := 0; i < len(expected); i++ {
			_, iexists := bullTracker[i]
			if !iexists {
				for j := 0; j < len(actual); j++ {
					_, jexists := bullTracker[j]
					if !jexists {
						if expected[i] == actual[j] {
							cows++
							break
						}
					}
				}
			}
		}
	}
	res := &GuessRes{cows, bulls}
	return res
}

//TODO
func VerifyWord(word string) bool {
	return true
}
