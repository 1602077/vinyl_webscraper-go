package records

import (
	"reflect"
	"testing"
)

var WKM = NewRecord("Tom Misch", "What Kinda Music", "", float32(30))
var LF = NewRecord("Jorja Smith", "Lost & Found", "", float32(100))
var NWBD = NewRecord("Loyle Carner", "Not Waving, But Drowning", "", float32(25))

func TestGetRecordMethods(t *testing.T) {
	artist, album, price := "tom misch", "geography", float32(25)
	rec := NewRecord(artist, album, "", price)

	if rec.GetArtist() != artist {
		t.Fatalf("GetArtist() = %s, Expected: %s", rec.GetArtist(), artist)
	}
	if rec.GetAlbum() != album {
		t.Fatalf("GetAlbum() = %s, Expected: %s", rec.GetAlbum(), album)
	}
	if rec.GetPrice() != price {
		t.Fatalf("GetPrice() = %v, Expected: %v", rec.GetPrice(), price)
	}
}

func TestSort(t *testing.T) {
	var records = Records{WKM, LF, NWBD}

	var tests = []struct {
		name       string
		sortMethod func(i *Record, j *Record) bool
		expected   Records
	}{
		{"ByArtist", ByArtist, Records{LF, NWBD, WKM}},
		{"ByAlbum", ByAlbum, Records{LF, NWBD, WKM}},
		{"ByPrice", ByPrice, Records{NWBD, WKM, LF}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sorted := records
			sorted.Sort(tt.sortMethod)

			if !reflect.DeepEqual(sorted, tt.expected) {
				t.Fatalf("%s failed\nExpected: %v, Got: %v", tt.name, tt.expected, sorted)
			}
		})
	}
}
