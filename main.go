package main

import (
	"encoding/json"
	"log"
	"time"
)

func main() {
	urls := readURLs("./data/input.txt")

	// Get current price of records in wishlist
	var r Records
	r = getRecords(urls)
	r.sortBy("artist")
	r.printRecords()

	bs, _ := json.MarshalIndent(r, "", " ")
	WriteToFile(bs, "./data/currentPrices.JSON")

	// Read in historical pricing and merge with current
	var rh RecordHistory
	rh, ReadErr := ReadFile("./data/allPrices.JSON", rh)
	if ReadErr != nil {
		log.Print("`./data/allPrices.JSON` does not exist; writing to new file")
	}

	today := time.Now().Format("2006-01-02")

	rh.MergeRecordHistories(RecordInstance{Date: today, Records: r})
	rh.sortBy("price")

	bs, _ = json.MarshalIndent(rh, "", " ")
	WriteToFile(bs, "./data/allPrices.JSON")
}
