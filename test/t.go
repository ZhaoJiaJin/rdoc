package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
    "net/url"
    "strings"
    "log"
)

func main() {

    apiurl := "http://localhost:8080/query"

	hc := http.Client{}
	//req, err := http.NewRequest("POST", APIURL, nil)

	form := url.Values{}
	form.Add("col", "Feeds")
	form.Add("q", `["all"]`)
	req, err := http.NewRequest("POST", apiurl, strings.NewReader(form.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, err := hc.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	fmt.Println(string(body))
}
