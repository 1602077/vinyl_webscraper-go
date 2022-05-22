// Data storage and handling methods.
package records

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"sort"
	"text/tabwriter"
)

type Record struct {
	artist      string
	album       string
	amazonUrl   string
	amazonPrice float32
}

type RecordJSON struct {
	Artist      string  `json:"artist"`
	Album       string  `json:"album"`
	AmazonUrl   string  `json:"amazon_url"`
	AmazonPrice float32 `json:"amazon_price"`
}

func (r *Record) MarshalJSON() ([]byte, error) {
	return json.Marshal(RecordJSON{
		r.artist,
		r.album,
		r.amazonUrl,
		r.amazonPrice,
	})
}

// TODO: Write Unit Test
func (r *Record) UnmarshalJSON(b []byte) error {
	tmp := &RecordJSON{}
	if err := json.Unmarshal(b, &tmp); err != nil {
		return err
	}
	r.artist = tmp.Artist
	r.album = tmp.Album
	r.amazonUrl = tmp.AmazonUrl
	r.amazonPrice = tmp.AmazonPrice
	return nil
}

func NewRecord(artist, album, url string, price float32) *Record {
	return &Record{
		artist:      artist,
		album:       album,
		amazonUrl:   url,
		amazonPrice: price,
	}
}

func (r *Record) GetArtist() string {
	return r.artist
}

func (r *Record) GetAlbum() string {
	return r.album
}

func (r *Record) GetPrice() float32 {
	return r.amazonPrice
}

type Records []*Record

func (r Records) MarshalJSON() ([]byte, error) {
	var rJson []*RecordJSON
	for _, rr := range r {
		rJson = append(rJson, &RecordJSON{
			Artist:      rr.artist,
			Album:       rr.album,
			AmazonUrl:   rr.amazonUrl,
			AmazonPrice: rr.amazonPrice,
		})
	}

	data, err := json.Marshal(rJson)
	if err != nil {
		return nil, err
	}

	return data, nil
}

// TODO: Write Unit Test
func (r Records) UnmarshalJSON(b []byte) error {
	var recordJsons []*RecordJSON
	if err := json.Unmarshal(b, &recordJsons); err != nil {
		log.Fatal("Records.UnmarshalJSON() failed: ", err)
	}

	for _, rr := range recordJsons {
		r = append(r, &Record{
			artist:      rr.Artist,
			album:       rr.Album,
			amazonUrl:   rr.AmazonUrl,
			amazonPrice: rr.AmazonPrice,
		})
	}

	return nil
}

func (r Records) PrintRecords() string {
	const format = "%v\t%v\t%v\n"

	var b bytes.Buffer
	tw := new(tabwriter.Writer).Init(&b, 0, 8, 4, ' ', 0)

	fmt.Fprintf(tw, format, "ARTIST", "ALBUM", "CURRENT PRICE")
	fmt.Fprintf(tw, format, "------", "-----", "-------------")
	for _, rr := range r {
		fmt.Fprintf(tw, format, rr.artist, rr.album, rr.amazonPrice)
	}
	tw.Flush()

	fmt.Println(b.String())
	return b.String()
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
