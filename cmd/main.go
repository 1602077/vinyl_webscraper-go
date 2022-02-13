package main

import (
	db "github.com/1602077/webscraper/pkg/postgres"
	r "github.com/1602077/webscraper/pkg/records"
	ws "github.com/1602077/webscraper/pkg/webscraper"
)

const (
	DBNAME = "prod"
)

func main() {
	urls := ws.ReadURLs("../data/input.txt")

	// Get current price of records in wishlist
	var currPrices r.Records
	currPrices = ws.GetRecords(urls)

	currPrices.Sort(r.ByArtist)
	currPrices.PrintRecords()

	pg := db.NewPostgresCli(DBNAME).Connect()
	defer pg.Close()

	for _, rec := range currPrices {
		pg.InserRecordAllTables(rec)
	}
}
