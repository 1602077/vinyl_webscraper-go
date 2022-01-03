package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/gocolly/colly"
)

type PageInfo struct {
	url, Artist, Album string
	// Price              float32
	Price string
}

var Pages []PageInfo

func main() {
	urls := readURLs("./input.txt")
	for _, u := range urls {
		p := getAmazonPageInfo(u)
		Pages = append(Pages, p)
	}

	j, _ := json.MarshalIndent(Pages, "", "	")
	// account for MarshalIndent escaping html
	j = bytes.Replace(j, []byte("\\u003c"), []byte("<"), -1)
	j = bytes.Replace(j, []byte("\\u003e"), []byte(">"), -1)
	j = bytes.Replace(j, []byte("\\u0026"), []byte("&"), -1)

	// write json byte slice to file
	f, err := os.Create("./output.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	n, err := f.Write(j)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("wrote %d bytes to \"./output.txt\"", n)
}

func getAmazonPageInfo(url string) (page PageInfo) {
	c := colly.NewCollector()

	c.OnHTML(`div[id=centerCol]`, func(e *colly.HTMLElement) {
		log.Println("visiting", url)
		album := e.ChildText(`span[id=productTitle]`)
		if album == "" {
			log.Println("no title found", e.Request.URL)
		}

		artist := e.ChildTexts(`a.a-link-normal`)[0]
		if artist == "" {
			log.Println("no artist found", e.Request.URL)
		}

		price := e.ChildText(`span[class='a-offscreen']`)
		if price == "" {
			log.Println("no price found", e.Request.URL)
		}

		page = PageInfo{
			url:    url,
			Album:  strings.Replace(album, " [VINYL]", "", 1),
			Price:  price,
			Artist: artist,
		}
	})
	c.Visit(url)
	return
}

func readURLs(filename string) []string {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}
	d := strings.Split(string(data), "\n")
	return d[:len(d)-1]
}

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
