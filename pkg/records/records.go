// Data storage and handling methods.
package records

import (
	"fmt"
	"os"
	"sort"
	"text/tabwriter"
)

type Record struct {
	artist, album string
	amazonUrl     string
	amazonPrice   float32
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

func (r Records) PrintRecords() {
	const format = "%v\t%v\t%v\n"
	tw := new(tabwriter.Writer).Init(os.Stdout, 0, 8, 4, ' ', 0)
	fmt.Println()
	fmt.Fprintf(tw, format, "ARTIST", "ALBUM", "CURRENT PRICE")
	fmt.Fprintf(tw, format, "------", "-----", "-------------")
	for _, rr := range r {
		fmt.Fprintf(tw, format, rr.artist, rr.album, rr.amazonPrice)
	}
	tw.Flush()
	fmt.Println()
	fmt.Println()
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
