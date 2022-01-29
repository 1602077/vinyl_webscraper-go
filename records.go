// Data storage and handling methods.
package main

import (
	"bytes"
	"encoding/json"
	"log"
	"os"
	"sort"
	"time"

	"github.com/go-gota/gota/dataframe"
)

type Record struct {
	Artist, Album          string
	amazonUrl, AmazonPrice string
}

type Records []Record

func (r Records) sortBy(field string) {
	sortRecordsByField(r, field)
}

func (r Records) writeToJSON(outname string) {
	j, _ := json.MarshalIndent(r, "", "	")
	cleanWriteJSON(j, outname)
}

func (r Records) ConvertToDataFrame() dataframe.DataFrame {
	return dataframe.LoadStructs(r)
}

type RecordInstance struct {
	Date    time.Time
	Records Records
}

type RecordHistory []RecordInstance

func (rh RecordHistory) writeToJSON(outname string) {
	j, _ := json.MarshalIndent(rh, "", " ")
	cleanWriteJSON(j, outname)
}

func (rh *RecordHistory) ReadFromJSON(filename string) (ReadErr error) {
	f, ReadErr := os.ReadFile(filename)
	if ReadErr != nil {
		return ReadErr
	}
	err := json.Unmarshal(f, &rh)
	if err != nil {
		log.Fatal(err)
	}
	return nil
}

func (rh1 *RecordHistory) MergeRecordHistories(ri RecordInstance) {
	*rh1 = append(*rh1, ri)
}

func (rh RecordHistory) sortBy(field string) {
	for _, v := range rh {
		r := v.Records
		sortRecordsByField(r, field)
	}
}

func sortRecordsByField(r Records, field string) {
	switch field {
	case "artist":
		sort.Slice(r, func(i, j int) bool {
			return r[i].Artist < r[j].Artist
		})
	case "album":
		sort.Slice(r, func(i, j int) bool {
			return r[i].Album < r[j].Album
		})
	case "price":
		sort.Slice(r, func(i, j int) bool {
			return r[i].AmazonPrice < r[j].AmazonPrice
		})
	default:
		sort.Slice(r, func(i, j int) bool {
			return r[i].Artist < r[j].Artist
		})
	}
}

// Takes a JSON object in a byte slice, escpaes all characters, and writes to file
func cleanWriteJSON(j []byte, outname string) {
	j = bytes.Replace(j, []byte("\\u003c"), []byte("<"), -1)
	j = bytes.Replace(j, []byte("\\u003e"), []byte(">"), -1)
	j = bytes.Replace(j, []byte("\\u0026"), []byte("&"), -1)

	f, err := os.Create(outname)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	n, err := f.Write(j)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("`%s` written (%v bytes)", outname, n)
}
