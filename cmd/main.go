package main

import (
	db "github.com/1602077/webscraper/pkg/postgres"
	r "github.com/1602077/webscraper/pkg/records"
	ws "github.com/1602077/webscraper/pkg/webscraper"
)

const DBNAME = "prod"

func main() {
	urls := ws.ReadURLs("../data/input.txt")

	var currPrices r.Records
	currPrices = ws.GetRecords(urls)

	pg := db.NewPostgresCli(DBNAME).Connect()
	// defer pg.Close()

	// type(curPrices): r.Records -> []*Record
	// type(rec): *r.Record)
	for _, rec := range currPrices {
		pg.InsertRecordAllTables(rec)
	}
	pg.PrintCurrentRecordPrices()
	pg.Close()
}
