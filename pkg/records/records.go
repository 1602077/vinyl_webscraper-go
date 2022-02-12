// Data storage and handling methods.
package records

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
	artist, album          string
	amazonUrl, amazonPrice string
}

func NewRecord(artist, album, url, price string) *Record {
	return &Record{
		artist:      artist,
		album:       album,
		amazonUrl:   url,
		amazonPrice: price,
	}
}

type Records []*Record

func (r Records) PrintRecords() {
	const format = "%v\t%v\t%v\n"
	tw := new(tabwriter.Writer).Init(os.Stdout, 0, 8, 4, ' ', 0)
	fmt.Fprintf(tw, format, "ARTIST", "ALBUM", "PRICE")
	fmt.Fprintf(tw, format, "------", "-----", "-----")
	for _, rr := range r {
		fmt.Fprintf(tw, format, rr.artist, rr.album, rr.amazonPrice)
	}
	tw.Flush()
}

type RecordsSort struct {
	r    []*Record
	less func(i, j *Record) bool
}

func (r RecordsSort) Len() int           { return len(r.r) }
func (r RecordsSort) Swap(i, j int)      { r.r[i], r.r[j] = r.r[j], r.r[i] }
func (r RecordsSort) Less(i, j int) bool { return r.less(r.r[i], r.r[j]) }

func ByArtist(i, j *Record) bool { return i.artist < j.artist }
func ByAlbum(i, j *Record) bool  { return i.album < j.album }
func ByPrice(i, j *Record) bool  { return i.amazonPrice < j.amazonPrice }

func (r Records) Sort(ByField func(*Record, *Record) bool) {
	sort.Sort(RecordsSort{r, ByField})
}

// Store of record wishlist data at a current instance in time.
type RecordInstance struct {
	Date    string
	Records Records
}

type RecordHistory []RecordInstance

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

func (rh RecordHistory) Sort(ByField func(*Record, *Record) bool) {
	for _, v := range rh {
		r := v.Records
		r.Sort(ByField)
	}
}

// PERSISTENCE FUNCTIONS

// Write byte slice to file, converting all chars back to utf-8
func WriteToFile(j []byte, filename string) {
	j = bytes.Replace(j, []byte("\\u003c"), []byte("<"), -1)
	j = bytes.Replace(j, []byte("\\u003e"), []byte(">"), -1)
	j = bytes.Replace(j, []byte("\\u0026"), []byte("&"), -1)

	f, err := os.Create(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	n, err := f.Write(j)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("`%s` written (%v bytes)", filename, n)
}

// Reads in JSON data to RecordHistory struct
func ReadFile(filename string, rh RecordHistory) (RecordHistory, error) {
	f, ReadErr := os.ReadFile(filename)
	if ReadErr != nil {
		return nil, ReadErr
	}
	err := json.Unmarshal(f, &rh)
	if err != nil {
		log.Fatal(err)
	}
	return rh, nil
}
