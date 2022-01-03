package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"sort"
	"strings"

	"github.com/gocolly/colly"
)

func main() {
	urls := readURLs("./data/input.txt")
	var p Records
	p = getRecords(urls)
	p.sortBy("artist")
	p.writeToJSON("./data/output.json")
}

type Record struct {
	Artist, Album          string
	amazonUrl, AmazonPrice string
}

type Records []Record

func (r Records) sortBy(field string) {
	switch field {
	case "artist":
		sort.Slice(r, func(i, j int) bool {
			return r[i].Artist < r[j].Artist
		})
	case "album":
		sort.Slice(r, func(i, j int) bool {
			return r[i].Album < r[j].Album
		})
	case "price":
		sort.Slice(r, func(i, j int) bool {
			return r[i].AmazonPrice < r[j].AmazonPrice
		})
	default:
		sort.Slice(r, func(i, j int) bool {
			return r[i].Artist < r[j].Artist
		})
	}
}

func (r Records) writeToJSON(outname string) {
	j, _ := json.MarshalIndent(r, "", "	")
	// account for MarshalIndent escaping html
	j = bytes.Replace(j, []byte("\\u003c"), []byte("<"), -1)
	j = bytes.Replace(j, []byte("\\u003e"), []byte(">"), -1)
	j = bytes.Replace(j, []byte("\\u0026"), []byte("&"), -1)

	// write json byte slice to file
	f, err := os.Create(outname)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	n, err := f.Write(j)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("wrote %d bytes to %s", n, outname)
}

func getRecords(urls []string) (records []Record) {
	ch := make(chan Record, len(urls))
	for _, u := range urls {
		go func(u string) {
			var r Record
			r = getAmazonPageInfo(u)
			ch <- r
		}(u)
	}
	for range urls {
		r := <-ch
		records = append(records, r)
	}
	return records
}

func getAmazonPageInfo(url string) (r Record) {
	c := colly.NewCollector()

	c.OnHTML(`div[id=centerCol]`, func(e *colly.HTMLElement) {
		// log.Println("visiting", url)
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

		r = Record{
			Album:       strings.Replace(album, " [VINYL]", "", 1),
			Artist:      artist,
			amazonUrl:   url,
			AmazonPrice: price,
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
