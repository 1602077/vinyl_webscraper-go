package main

import (
	db "github.com/1602077/webscraper/pkg/postgres"
	r "github.com/1602077/webscraper/pkg/records"
	ws "github.com/1602077/webscraper/pkg/webscraper"
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