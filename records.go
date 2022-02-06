// Data storage and handling methods.
package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sort"
	"text/tabwriter"
)

type Record struct {
	Artist, Album          string
	amazonUrl, AmazonPrice string
}

type Records []*Record

func (r Records) writeToJSON(outname string) {
	j, _ := json.MarshalIndent(r, "", "	")
	cleanWriteJSON(j, outname)
}

// Records sort.Interfaces
type byArtist []*Record

func (x byArtist) Len() int           { return len(x) }
func (x byArtist) Less(i, j int) bool { return x[i].Artist < x[j].Artist }
func (x byArtist) Swap(i, j int)      { x[i], x[j] = x[j], x[i] }

type byAlbum []*Record

func (x byAlbum) Len() int           { return len(x) }
func (x byAlbum) Less(i, j int) bool { return x[i].Album < x[j].Album }
func (x byAlbum) Swap(i, j int)      { x[i], x[j] = x[j], x[i] }

type byAmazonPrice []*Record

func (x byAmazonPrice) Len() int           { return len(x) }
func (x byAmazonPrice) Less(i, j int) bool { return x[i].AmazonPrice < x[j].AmazonPrice }
func (x byAmazonPrice) Swap(i, j int)      { x[i], x[j] = x[j], x[i] }

func (r Records) printRecords() {
	const format = "%v\t%v\t%v\n"
	tw := new(tabwriter.Writer).Init(os.Stdout, 0, 8, 4, ' ', 0)
	fmt.Fprintf(tw, format, "ARTIST", "ALBUM", "PRICE")
	fmt.Fprintf(tw, format, "------", "-----", "-----")
	for _, rr := range r {
		fmt.Fprintf(tw, format, rr.Artist, rr.Album, rr.AmazonPrice)
	}
	tw.Flush()
}

// Store of record wishlist data at a current instance in time.
type RecordInstance struct {
	Date    string
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

func (rh *RecordHistory) MergeRecordHistories(ri RecordInstance) {
	var mergedRH RecordHistory
	currDate := ri.Date
	for _, r := range *rh {
		// Filter out recordInstance for date already exists
		if r.Date == currDate {
			continue
		}
		mergedRH = append(mergedRH, r)
	}
	mergedRH = append(mergedRH, ri)
	*rh = mergedRH
}

func (rh RecordHistory) sortBy(field string) {
	for _, v := range rh {
		r := v.Records
		switch field {
		case "artist":
			sort.Sort(byArtist(r))
		case "album":
			sort.Sort(byAlbum(r))
		case "price":
			sort.Sort(byAmazonPrice(r))
		default:
			sort.Sort(byArtist(r))
		}
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
