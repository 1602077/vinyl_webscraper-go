package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"time"
)

const amazonBaseURL = "https://www.amazon.co.uk/s"

func createURL(baseURL, recordName string) string {
	var u *url.URL
	var err error
	u, err = url.Parse(baseURL)
	if err != nil {
		log.Fatal(err)
	}
	params := url.Values{}
	params.Add("k", recordName+" vinyl")
	u.RawQuery = params.Encode()
	return u.String()
}

func getHTML(url string) string {
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := client.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	return string(body)
}
