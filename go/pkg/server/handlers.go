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
	flag.StringVar(&ENV_FILEPATH, "env", "../.env", "sets environment config (.env) filepath")
	flag.Parse()
	fmt.Printf("runtime config filepath: '%s'\n", ENV_FILEPATH)
}

func RefreshRecordPrices(w http.ResponseWriter, r *http.Request) {
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
