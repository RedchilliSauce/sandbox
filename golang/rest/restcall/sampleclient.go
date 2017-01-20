package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
)

func main() {
	url := "http://localhost:12345/people"
	fmt.Println("Url : " + url)

	idURLAppend := "/3"
	var jsonStr = []byte(`{"firstname":"Jon", "lastname":"Snow", "address":{"city":"Castle Black", "state":"Wall"}}`)
	req, _ := http.NewRequest("POST", url+idURLAppend, bytes.NewBuffer(jsonStr))
	//req.Header.Set("X-Custom-Header", "myvalue")
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	fmt.Println("response Status:", resp.Status)
	fmt.Println("response Headers:", resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("response Body:", string(body))
}
