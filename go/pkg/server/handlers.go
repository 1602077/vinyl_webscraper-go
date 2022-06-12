// server packages api routing and handling for go webscraping app.
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

// GetRecords queries the Record information and their current prices for all
// records currently in the postgres database.
func GetRecords(w http.ResponseWriter, r *http.Request) {
	pg := postgres.GetPgInstance().Connect(ENV_FILEPATH)
	defer pg.Close()

	recs := pg.GetCurrentRecordPrices()

	/* HTML Rendered Site
	t, err := template.ParseFiles("../templates/records.html")
	if err != nil {
		log.Fatal("HomePage handler: ", err)
	}
	t.Execute(w, recs)
	*/

	recsJson, err := recs.MarshalJSON()
	if err != nil {
		log.Printf("err: HomePage handler: %s\n", err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(recsJson)
}

// PutRecords gets the current prices for all records in database, by
// making a calling to webscaper.GetRecords. All prices are written back to
// database and the record price information written to the http body.
func PutRecords(w http.ResponseWriter, r *http.Request) {
	var currPrices records.Records
	urls := webscraper.ReadURLs("../../input.txt")
	currPrices = webscraper.GetRecords(urls)

	pg := postgres.GetPgInstance().Connect(ENV_FILEPATH)
	defer pg.Close()
	for _, rec := range currPrices {
		pg.InsertRecord(rec)
	}
	pg.PrintCurrentPrices()

	cpJson, err := currPrices.MarshalJSON()
	if err != nil {
		log.Printf("err: GetRecordPrices handler: %s\n", err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(cpJson)
}

// GetRecord takes an input record id and returns the record information (i.e.
// artist, album) and it's full pricing history.
func GetRecord(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	urlVars := mux.Vars(r)
	rId, err := strconv.Atoi(urlVars["id"])
	if err != nil {
		log.Printf("GetRecord Handler: %s\n", err)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	pg := postgres.GetPgInstance().Connect(ENV_FILEPATH)
	defer pg.Close()

	var rph *records.RecordPriceHistory
	rph = pg.GetRecordPriceHistory(rId)
	if rph == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	rphJson, err := json.Marshal(rph)
	if err != nil {
		log.Printf("err: GetRecord: %s\n", err)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(rphJson)
}
