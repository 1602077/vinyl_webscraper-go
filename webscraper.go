// Webscraper to scrape record information from an amazon URL.
package main

import (
	"io/ioutil"
	"log"
	"regexp"
	"strings"

	"github.com/gocolly/colly"
)

// Reads from filename a list of urls each seperated by a `\n`.
func readURLs(filename string) []string {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}
	d := strings.Split(string(data), "\n")
	return d[:len(d)-1]
}

// Gets the Artist, Album Name and Price for a given record from amazon URL.
func getAmazonPageInfo(url string) (r *Record) {
	c := colly.NewCollector()

	c.OnHTML(`div[id=centerCol]`, func(e *colly.HTMLElement) {
		album := e.ChildText(`span[id=productTitle]`)
		if album == "" {
			log.Println("no title found", e.Request.URL)
		}

		artist := e.ChildText(`a.a-link-normal`)
		if artist == "" {
			log.Println("no artist found", e.Request.URL)
		}

		price := e.ChildText(`span[class='a-offscreen']`)
		if price == "" {
			log.Println("no price found", e.Request.URL)
		}

		r = &Record{
			Album:       strings.Replace(album, " [VINYL]", "", 1),
			Artist:      parseArtist(artist),
			amazonUrl:   url,
			AmazonPrice: parsePrice(price),
		}
	})
	c.Visit(url)
	return
}

func parseArtist(s string) string {
	re := regexp.MustCompile(` \d+ ratings`)
	indx := re.FindStringIndex(s)[0]
	return s[:indx]
}

func parsePrice(s string) string {
	re := regexp.MustCompile(`£[\d.]+`)
	return re.FindString(s)
}

// Concurrently calls `getAmazonPageInfo` for a list of URLS.
func getRecords(urls []string) (records Records) {
	// limit to 10 concurrent requests at a time.
	// ch := make(chan *Record, len(urls))
	ch := make(chan *Record, 10)
	for _, u := range urls {
		go func(u string) {
			var r *Record
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
