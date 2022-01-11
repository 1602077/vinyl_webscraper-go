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
	r.sortBy("artist")
	r.writeToJSON("./data/currentPrices.JSON")

	// Append current and historical pricing
	var rh RecordHistory
	if ReadErr := rh.ReadFromJSON("./data/allPrices.JSON"); ReadErr != nil {
		log.Print("`./data/allPrices.JSON` does not exist; writing to new file")
	}
	rh.MergeRecordHistories(RecordInstance{Date: time.Now(), Records: r})
	rh.sortBy("artist")
	rh.writeToJSON("./data/allPrices.JSON")

	// TODO: Pretty print output into a neat table comparing historical prices
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

func (r Records) sortBy(field string) {
	sortRecordsByField(r, field)
}

func (r Records) writeToJSON(outname string) {
	j, _ := json.MarshalIndent(r, "", "	")
	cleanseWriteJSON(j, outname)
}

type RecordInstance struct {
	Date    time.Time
	Records Records
}

type RecordHistory []RecordInstance

func (rh RecordHistory) writeToJSON(outname string) {
	j, _ := json.MarshalIndent(rh, "", " ")
	cleanseWriteJSON(j, outname)
}

func (rh *RecordHistory) ReadFromJSON(filename string) (ReadErr error) {
	f, ReadErr := os.ReadFile(filename)
	if ReadErr != nil {
		return ReadErr
	}
	err := json.Unmarshal(f, &rh)
	if err != nil {
		log.Fatal(err)
	}
	return nil
}

func (rh1 *RecordHistory) MergeRecordHistories(ri RecordInstance) {
	*rh1 = append(*rh1, ri)
}

func (rh RecordHistory) sortBy(field string) {
	for _, v := range rh {
		r := v.Records
		sortRecordsByField(r, field)
	}
}

func sortRecordsByField(r Records, field string) {
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

// Takes a JSON object in a byte slice, escpaes all characters, and writes to file
func cleanseWriteJSON(j []byte, outname string) {
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
	log.Printf("`%s` written (%v bytes)", outname, n)
}
