package records

import (
	"bytes"
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

func TestPrintRecords(t *testing.T) {
	records := Records{
		WKM,
		LF,
		NWBD,
	}
	actual := records.Print()

	if len(actual) != 258 {
		t.Fatalf("Expected length of 258 for PrintRecords(), Got: %v", len(actual))
	}
}

func TestRecordMarshalJSON(t *testing.T) {
	original := WKM

	marshalled, err := original.MarshalJSON()
	if err != nil {
		t.Fatal("Record.MarshalJSON() returned an error, when non was expected.")
	}

	t.Run("Record.MarshalJSON()", func(t *testing.T) {
		expected := []byte(`{"artist":"Tom Misch","album":"What Kinda Music","amazon_url":"","amazon_price":30}`)
		res := bytes.Compare(marshalled, expected)
		if res != 0 {
			t.Fatalf("Expected: %v\nGot: %v\n", expected, marshalled)
		}
	})

	t.Run("Record.UnmarshalJSON()", func(t *testing.T) {
		unmarshalled := &Record{}
		unmarshalled.UnmarshalJSON(marshalled)
		if !reflect.DeepEqual(original, unmarshalled) {
			t.Fatalf("Expected: %v\nGot: %v\n", original, unmarshalled)
		}
	})

}

func TestRecordsMarshalJSON(t *testing.T) {
	original := Records{
		WKM,
		LF,
	}

	marshalled, err := original.MarshalJSON()
	if err != nil {
		t.Fatal("Records.MarshalJSON() returned an error, when non was expected.")
	}

	t.Run("Records.MarshalJSON()", func(t *testing.T) {
		expected := []byte(`[{"artist":"Tom Misch","album":"What Kinda Music","amazon_url":"","amazon_price":30},{"artist":"Jorja Smith","album":"Lost \u0026 Found","amazon_url":"","amazon_price":100}]`)
		res := bytes.Compare(marshalled, expected)
		if res != 0 {
			t.Fatalf("Expected: %v\nGot: %v\n", expected, marshalled)
		}
	})

	t.Run("Records.UnmarshalJSON()", func(t *testing.T) {
		unmarshalled := make(Records, 0)
		unmarshalled.UnmarshalJSON(marshalled)

		if !reflect.DeepEqual(original, unmarshalled) {
			t.Fatalf("Expected: %v\nGot: %v\n", original, unmarshalled)
		}
	})

}
