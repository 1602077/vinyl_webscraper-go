package server

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"text/template"

	"github.com/1602077/webscraper/go/pkg/postgres"
	"github.com/1602077/webscraper/go/pkg/records"
	"github.com/1602077/webscraper/go/pkg/webscraper"
)

var ENV_FILEPATH string

func init() {
	flag.StringVar(&ENV_FILEPATH, "env", "../../.env", "sets environment config (.env) filepath")
	flag.Parse()
	fmt.Printf("runtime config filepath: '%s'\n", ENV_FILEPATH)
}

func HomePage(w http.ResponseWriter, r *http.Request) {
	pg := postgres.GetPgInstance().Connect(ENV_FILEPATH)
	defer pg.Close()

	var recs records.Records
	recs = postgres.ReadQueryToRecords(pg.GetCurrentRecordPrices())

	t, err := template.ParseFiles("../templates/records.html")
	if err != nil {
		log.Fatal("HomePage handler: ", err)
	}
	t.Execute(w, recs)
}

func GetRecordPrices(w http.ResponseWriter, r *http.Request) {
	// wd := postgres.GetEnVar(ENV_FILEPATH, "WORKDIR")
	// Get record price data
	var currPrices records.Records
	urls := webscraper.ReadURLs("../../input.txt")
	currPrices = webscraper.GetRecords(urls)

	// Write to postgres
	pg := postgres.GetPgInstance().Connect(ENV_FILEPATH)
	defer pg.Close()
	for _, rec := range currPrices {
		pg.InsertRecord(rec)
	}
	pg.PrintCurrentRecordPrices()

	cpJson, err := currPrices.MarshalJSON()
	if err != nil {
		log.Printf("GetRecordPrices handler: %s\n", err)
	}

	// Write header
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(cpJson)
}
