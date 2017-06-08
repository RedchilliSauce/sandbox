package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/labstack/echo"
)

//Game ...
type Game struct {
	Player1 string              `json:"p1"`
	Player2 string              `json:"p2"`
	Word    string              `json:"word"`
	Guesses map[string]GuessRes `json:"guesses"`
}

//GuessRes ...
type GuessRes struct {
	Cows  int64 `json:"cows"`
	Bulls int64 `json:"bulls"`
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

	Games = make(map[string]Game)
	e := echo.New()

	e.GET("/", index)
	e.GET("/newgame", newgame)
	e.GET("/existinggame", existinggame)
	e.GET("/guess/:p1/:p2", guess)
	e.POST("/getguessres", getguessres)
	e.POST("/creategame", creategame)

	e.Start(":" + port)
}

func index(c echo.Context) error {
	html := `
	<!DOCTYPE html>
	<html>
	<h1>Cows n Bulls</h1>

	<body>
		<a href="/newgame">New Game</a>
		<br>
		<a href="/existinggame">Continue existing game</a>
	</body>

	</html>
	`
	return c.HTML(http.StatusOK, html)
}

func newgame(c echo.Context) error {
	html := `<!DOCTYPE html>
	<html>
    <body>
        <form action="/creategame" method="post">
            Player1<input type="text" required name="p1">
            <br>Player2
            <input type="text" required name="p2">
            <br>Word to be guessed
            <input type="text" required name="word">
            <br><br>
            <input type="submit" value="Submit">
        </form>
        <br><br>
        <a href="/">Back to Main Menu</a>
    </body>
    </html>`

	return c.HTML(http.StatusOK, html)
}

func existinggame(c echo.Context) error {
	rowtmpl := `<tr><td>%s</td><td>%s</td><td><a href="%s">Link</td></tr><br>`
	headertmpl := `<tr><th>%s</th><th>%s</th><th>%s</th></tr><br>`

	html := `<!DOCTYPE html>
	<html>
	<h1>Cows n Bulls</h1>

	<body>
		<table>
		<style>
	table {
    font-family: arial, sans-serif;
    border-collapse: collapse;
    width: 100%;
}

td, th {
    border: 1px solid #dddddd;
    text-align: left;
    padding: 8px;
}
	</style>
	` + fmt.Sprintf(headertmpl, "Player1", "Player2", "Click to play")

	for _, game := range Games {
		link := `/guess/` + game.Player1 + `/` + game.Player2
		row := fmt.Sprintf(rowtmpl, game.Player1, game.Player2, link)
		html = html + row
	}

	html = html + `</table>
	<br><br>
    <a href="/">Back to Main Menu</a>
	</body>

	</html>`
	return c.HTML(http.StatusOK, html)
}

func guess(c echo.Context) error {
	p1 := c.Param("p1")
	p2 := c.Param("p2")
	key := generateKey(p1, p2)

	game, _ := Games[key]
	var html string
	html = `<!DOCTYPE html>
	<html>
	<body>
    <form action="/getguessres" method="post">
        Your guess<br>
        <input type="text" required name="guess">
        <br><br>
		<input type="hidden" name="p1" value="` + p1 + `">
		<input type="hidden" name="p2" value="` + p2 + `">
        <input type="submit" value="Submit">
    </form>
	` + getResultOutputAppender(game) + `
	<br><br>
    <a href="/">Back to Main Menu</a>
</body>
</html>`

	return c.HTML(http.StatusOK, html)
}

func getguessres(c echo.Context) error {
	p1 := c.FormValue("p1")
	p2 := c.FormValue("p2")
	guess := c.FormValue("guess")

	key := generateKey(p1, p2)

	game, exists := Games[key]

	if exists {
		res := getGuessRes(game.Word, guess)
		if res == nil {
			return c.String(http.StatusBadRequest, "Incorrect input")
		}
		if res.Bulls == 4 && res.Cows == 0 {
			delete(Games, key)
			return c.HTML(http.StatusOK, `<!DOCTYPE html>
	<html>
	<body>
	Whee you guessed the word correctly. Word was *`+game.Word+`*
	<br><br>
    <a href="/">Back to Main Menu</a>
	</body></html>`)
		}
		game.Guesses[guess] = *res
	}

	html := `<!DOCTYPE html>
	<html>
	<body>
	<form action="/getguessres" method="post">
        Your guess<br>
        <input type="text" required name="guess">
        <br><br>
		<input type="hidden" name="p1" value="` + p1 + `">
		<input type="hidden" name="p2" value="` + p2 + `">
        <input type="submit" value="Submit">
    </form>
	` + getResultOutputAppender(game) + `
	<br><br>
    <a href="/">Back to Main Menu</a>
		</body></html>
	`
	return c.HTML(http.StatusOK, html)
}

func generateKey(p1 string, p2 string) string {
	return p1 + "#@#" + p2
}

func creategame(c echo.Context) error {
	guesses := make(map[string]GuessRes)

	p1 := c.FormValue("p1")
	p2 := c.FormValue("p2")
	word := c.FormValue("word")

	key := generateKey(p1, p2)
	newGame := Game{p1, p2, word, guesses}

	_, exists := Games[key]
	if exists {
		return c.HTML(http.StatusBadRequest, `<html>
    <body>
        Game already exists between the two players
        <br><br>
        <a href="/">Back to Main Menu</a>
    </body>
    </html>`)
	}

	Games[key] = newGame

	html := `<html>
    <body>
        Game successfully created
        <br><br>
        <a href="/">Back to Main Menu</a>
    </body>
    </html>`
	return c.HTML(http.StatusOK, html)
}

func getGuessRes(expected string, actual string) *GuessRes {
	if len(expected) != len(actual) {
		return nil
	}

	cows := int64(0)
	bulls := int64(0)

	expected = strings.ToUpper(expected)
	actual = strings.ToUpper(actual)

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

//VerifyWord TODO
func VerifyWord(word string) bool {
	return true
}

func getResultOutputAppender(game Game) string {
	rowtmpl := `<tr><td>%s</td><td>%s</td><td>%s</td></tr><br>`
	headertmpl := `<tr><th>%s</th><th>%s</th><th>%s</th></tr><br>`

	html := `
	<style>
	table {
    font-family: arial, sans-serif;
    border-collapse: collapse;
    width: 100%;
}

td, th {
    border: 1px solid #dddddd;
    text-align: left;
    padding: 8px;
}
	</style>
	<table>
	` + fmt.Sprintf(headertmpl, "Word", "Cows", "Bulls")

	for word, guess := range game.Guesses {
		row := fmt.Sprintf(rowtmpl, word, strconv.FormatInt(guess.Cows, 10), strconv.FormatInt(guess.Bulls, 10))
		html = html + row
	}

	html = html + `</table>
	`

	return html
}
