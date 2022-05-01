package main

import (
	db "github.com/1602077/webscraper/go/pkg/postgres"
	r "github.com/1602077/webscraper/go/pkg/records"
	ws "github.com/1602077/webscraper/go/pkg/webscraper"
)

const ENV_FILEPATH = "../.env"

func main() {
	wd := db.GetEnVar(ENV_FILEPATH, "WORKDIR")
	urls := ws.ReadURLs(wd + "/input.txt")

	var currPrices r.Records
	currPrices = ws.GetRecords(urls)

	pg := db.GetPgInstance().Connect(ENV_FILEPATH)
	for _, rec := range currPrices {
		pg.InsertRecord(rec)
	}
	pg.PrintCurrentRecordPrices()
	pg.Close()
}
