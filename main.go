package main

import (
	"fmt"
	"log"
	"time"
)

func main() {
	urls := readURLs("./data/input.txt")

	// Get current price of records in wishlist
	var r Records
	r = getRecords(urls)
	r.sortBy("artist")
	r.writeToJSON("./data/currentPrices.JSON")

	fmt.Printf("\nCURRENT PRICE DATAFRAME:\n%v\n", r.ConvertToDataFrame())

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
