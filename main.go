package main

import (
	"log"
	"sort"
	"time"
)

func main() {
	urls := readURLs("./data/input.txt")

	// Get current price of records in wishlist
	var r Records
	r = getRecords(urls)
	sort.Sort(byArtist(r))
	r.writeToJSON("./data/currentPrices.JSON")
	r.printRecords()

	// Append current and historical pricing
	var rh RecordHistory
	if ReadErr := rh.ReadFromJSON("./data/allPrices.JSON"); ReadErr != nil {
		log.Print("`./data/allPrices.JSON` does not exist; writing to new file")
	}

	today := time.Now().Format("2006-01-02")

	rh.MergeRecordHistories(RecordInstance{Date: today, Records: r})
	rh.sortBy("price")
	rh.writeToJSON("./data/allPrices.JSON")
}
