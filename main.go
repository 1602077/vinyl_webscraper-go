package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/gocolly/colly"
)

func main() {
	urls := readURLs("./data/input.txt")
	// Get current price of records in wishlist
	var r Records
	r = getRecords(urls)
	r.writeToJSON("./data/currentPrices.JSON")

	// Append current and historical pricing
	// TODO: Reimplement sorting function for `RecordHistory` Type
	var rh RecordHistory
	rh.ReadFromJSON("./data/output.json")
	rh.MergeRecordHistories(RecordInstance{Date: time.Now(), Records: r})
	rh.writeToJSON("./data/allPrices.JSON")
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

type Record struct {
	Artist, Album          string
	amazonUrl, AmazonPrice string
}

type Records []Record

// TODO: Reimplement to account for nested JSON with date
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
	log.Printf("`Records` written to %s (%d bytes)", outname, n)
}

type RecordInstance struct {
	Date    time.Time
	Records Records
}

type RecordHistory []RecordInstance

// TODO: Look into refactoring this & Record method
func (rh RecordHistory) writeToJSON(outname string) {
	j, _ := json.MarshalIndent(rh, "", " ")

	j = bytes.Replace(j, []byte("\\u003c"), []byte("<"), -1)
	j = bytes.Replace(j, []byte("\\u003e"), []byte(">"), -1)
	j = bytes.Replace(j, []byte("\\u0026"), []byte("&"), -1)

	f, err := os.Create(outname)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	n, err := f.Write(j)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("`RecordHistory` written to %s (%d bytes)", outname, n)
}

// Reads in historical record pricing data from a saved JSON back into &rh
func (rh *RecordHistory) ReadFromJSON(filename string) {
	f, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}
	err = json.Unmarshal(f, &rh)
	if err != nil {
		log.Fatal(err)
	}
}

func (rh1 *RecordHistory) MergeRecordHistories(ri RecordInstance) {
	*rh1 = append(*rh1, ri)
}
