package server

import (
	"log"
	"net/http"
	"text/template"

	"github.com/1602077/webscraper/go/pkg/postgres"
	"github.com/1602077/webscraper/go/pkg/records"
	"github.com/1602077/webscraper/go/pkg/webscraper"
)

const ENV_FILEPATH = "../.env"

func RefreshRecordPrices(w http.ResponseWriter, r *http.Request) {
	// Get current prices of all records in input.txt
	wd := postgres.GetEnVar(ENV_FILEPATH, "WORKDIR")
	urls := webscraper.ReadURLs(wd + "/input.txt")

	var currPrices records.Records
	currPrices = webscraper.GetRecords(urls)

	pg := postgres.GetPgInstance().Connect(ENV_FILEPATH)
	defer pg.Close()

	for _, rec := range currPrices {
		pg.InsertRecord(rec)
	}
	pg.PrintCurrentRecordPrices()

	t, err := template.ParseFiles("templates/records.html")
	if err != nil {
		log.Fatal("RefreshRecordPrices: ", err)
	}
	t.Execute(w, currPrices)
}

func DisplayRecordPrices(w http.ResponseWriter, r *http.Request) {
	pg := postgres.GetPgInstance().Connect(ENV_FILEPATH)
	defer pg.Close()

	var recs records.Records
	recs = postgres.ReadQueryToRecords(pg.GetCurrentRecordPrices())

	t, err := template.ParseFiles("templates/records.html")
	if err != nil {
		log.Fatal("DisplayRecordPrices: ", err)
	}
	t.Execute(w, recs)
}
