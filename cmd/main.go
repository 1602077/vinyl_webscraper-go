package main

import (
	"encoding/json"
	"log"
	"time"

	r "github.com/1602077/webscraper/pkg/records"
	ws "github.com/1602077/webscraper/pkg/webscraper"
)

func main() {
	urls := ws.ReadURLs("../data/input.txt")

	// Get current price of records in wishlist
	var currPrices r.Records
	currPrices = ws.GetRecords(urls)
	currPrices.Sort(r.ByArtist)
	currPrices.PrintRecords()

	bs, _ := json.MarshalIndent(currPrices, "", " ")
	r.WriteToFile(bs, "../data/currentPrices.JSON")

	// Read in historical pricing and merge with current
	var histPrices r.RecordHistory
	histPrices, ReadErr := r.ReadFile("../data/allPrices.JSON", histPrices)
	if ReadErr != nil {
		log.Print("`../data/allPrices.JSON` does not exist; writing to new file")
	}

	today := time.Now().Format("2006-01-02")

	histPrices.MergeRecordHistories(r.RecordInstance{Date: today, Records: currPrices})
	histPrices.Sort(r.ByArtist)

	bs, _ = json.MarshalIndent(histPrices, "", " ")
	r.WriteToFile(bs, "../data/allPrices.JSON")
}
