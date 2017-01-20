package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

const omdbURL string = "http://www.omdbapi.com/?"

func search(title string) {
	finalURL := omdbURL + "t=" + title

	req, err := http.NewRequest("GET", finalURL, nil)

	if err != nil {
		panic(err)
	}

	client := &http.Client{}
	resp, er := client.Do(req)

	if er != nil {
		panic(er)
	}

	defer resp.Body.Close()

	fmt.Println("response status:", resp.Status)
	fmt.Println("response Headers:", resp.Header)

	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("response Body:", string(body))
}

func main() {
	var title string
	title = "Die Hard"

	titleFormatted := url.QueryEscape(title)
	fmt.Println(titleFormatted)
	search(titleFormatted)
}
