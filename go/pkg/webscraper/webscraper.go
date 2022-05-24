// Webscraper to scrape record information from an amazon URL.
package webscraper

import (
	"io/ioutil"
	"log"
	"regexp"
	"strconv"
	"strings"

	r "github.com/1602077/webscraper/go/pkg/records"
	"github.com/gocolly/colly"
)

// ReadURLs reads in  a list of urls each seperated by a `\n` from the input
// file to a slice of strings.
func ReadURLs(filename string) []string {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}
	d := strings.Split(string(data), "\n")
	return d[:len(d)-1]
}

// getAmazonPageInfo gets the Artist, Album Name and Price for a given record
// from an amazon URL by using the gocolly package.
func getAmazonPageInfo(url string) (pageinfo *r.Record) {
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

		pageinfo = r.NewRecord(
			parseArtist(artist),
			strings.Replace(album, " [VINYL]", "", 1),
			url,
			parsePrice(price),
		)
	})
	c.Visit(url)

	var emptyRecord *r.Record
	if emptyRecord == pageinfo {
		log.Fatal("`getAmazonPageInfo()` returned nil for all fields. Exceed call limit for session")
	}

	return
}

// parseArtist does a regex parse of the getAmazonPageInfo artist field output
// to remove the ratings tag which is occasionally included in html element.
func parseArtist(s string) string {
	re := regexp.MustCompile(` \d+,?\d+ ratings`)
	indx := re.FindStringIndex(s)[0]
	return s[:indx]
}

// parsePrice does a regex parse of the getAmazonPageInfo price to strip out
// any redudant text that may be lingering in the html element.
func parsePrice(s string) float32 {
	re := regexp.MustCompile(`[\d.]+`)
	price_str := re.FindString(s)
	flt, _ := strconv.ParseFloat(price_str, 32)
	return float32(flt)
}

// GetRecords concurrently calls getAmazonPageInfo to allow for the scraping of
// URLS to be performed in parallel.
func GetRecords(urls []string) (rs r.Records) {
	// limit to 10 concurrent requests at a time.
	ch := make(chan *r.Record, 10)
	for _, u := range urls {
		go func(u string) {
			var r *r.Record
			r = getAmazonPageInfo(u)
			ch <- r
		}(u)
	}
	for range urls {
		r := <-ch
		rs = append(rs, r)
	}
	return rs
}
