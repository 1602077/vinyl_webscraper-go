package main

import (
	"log"
	"os"

	db "github.com/1602077/webscraper/pkg/postgres"
	r "github.com/1602077/webscraper/pkg/records"
	ws "github.com/1602077/webscraper/pkg/webscraper"
)

const ENV_FILEPATH = ".env"

func main() {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	urls := ws.ReadURLs(dir + "/data/input.txt")

	var currPrices r.Records
	currPrices = ws.GetRecords(urls)

	pg := db.NewPostgresCli(ENV_FILEPATH).Connect()
	for _, rec := range currPrices {
		pg.InsertRecord(rec)
	}
	pg.PrintCurrentRecordPrices()
	pg.Close()
}
