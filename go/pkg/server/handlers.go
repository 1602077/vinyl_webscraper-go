package server

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/1602077/webscraper/go/pkg/postgres"
	"github.com/1602077/webscraper/go/pkg/records"
	"github.com/1602077/webscraper/go/pkg/webscraper"
	"github.com/gorilla/mux"
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

	/* HTML Rendered Site
	t, err := template.ParseFiles("../templates/records.html")
	if err != nil {
		log.Fatal("HomePage handler: ", err)
	}
	t.Execute(w, recs)
	*/

	// Create response body & header
	recsJson, err := recs.MarshalJSON()
	if err != nil {
		log.Printf("err: HomePage handler: %s\n", err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(recsJson)

}

func GetRecordPrices(w http.ResponseWriter, r *http.Request) {
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

	// Create response body & header
	cpJson, err := currPrices.MarshalJSON()
	if err != nil {
		log.Printf("err: GetRecordPrices handler: %s\n", err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(cpJson)
}

func GetRecord(w http.ResponseWriter, r *http.Request) {
	urlVars := mux.Vars(r)
	rId, err := strconv.Atoi(urlVars["id"])
	if err != nil {
		log.Printf("GetRecord Handler: %s\n", err)
	}

	var rph *records.RecordPriceHistory

	pg := postgres.GetPgInstance().Connect(ENV_FILEPATH)
	defer pg.Close()
	rph = pg.GetRecordPriceHistory(rId)

	rphJson, err := json.Marshal(rph)
	if err != nil {
		log.Printf("err: GetRecord: %s\n", err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(rphJson)
}
