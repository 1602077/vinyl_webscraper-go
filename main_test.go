package main

import (
	"fmt"
	"os"
	"testing"
	"time"
)

var WKM = Record{
	Artist:      "Tom Misch",
	Album:       "What Kinda Music",
	AmazonPrice: "£30",
}

var LF = Record{
	Artist:      "Jorja Smith",
	Album:       "Lost & Found",
	AmazonPrice: "£22.75",
}

var NWBD = Record{
	Artist:      "Loyle Carner",
	Album:       "Not Waving, But Drowning",
	AmazonPrice: "£25",
}

// Tests that embedded history records with date is formatted correctly
func TestRecordHistory(t *testing.T) {
	today := time.Now()
	records := Records{WKM, LF}

	data := RecordHistory{
		{Date: today, Records: records},
	}
	fmt.Println(data)
	data.writeToJSON("./data/TestRecordHistory.JSON")
}

func TestReadFROMJSON(t *testing.T) {
	filename := "./data/TestRecordHistory.JSON"
	var rh RecordHistory
	rh.ReadFromJSON(filename)

	if len(rh) == 0 {
		t.Errorf("Read of %s failed, length of record history is 0", filename)
	}
}

// Test that record read in from JSON can be merged current scrape
func TestMergeRecordHistories(t *testing.T) {
	filename := "./data/TestRecordHistory.JSON"
	var rh RecordHistory
	rh.ReadFromJSON(filename)

	today := time.Now()

	rh.MergeRecordHistories(RecordInstance{today, Records{NWBD}})
	rh.writeToJSON("./data/TestRecordHistoryMerge.JSON")

	if len(rh) != 2 {
		t.Errorf("Expected record history to contain  2 dates of scraping, got %v", len(rh))
	}
	os.Remove("./data/TestRecordHistory.JSON")
	os.Remove("./data/TestRecordHistoryMerge.JSON")
}

/*
func TestRecordHistorySortBy(t *testing.T) {
	var rh RecordHistory
	rh.ReadFromJSON("./data/allPrices.JSON")

	rh.sortBy("album")
	fmt.Print(rh)
}

type urlTest struct {
	baseURL, recordName, expected string
}

var urlTests = []urlTest{
	{amazonBaseURL, "what kinda music", "https://www.amazon.co.uk/s?k=what+kinda+music+vinyl"},
	{amazonBaseURL, "venice", "https://www.amazon.co.uk/s?k=venice+vinyl"},
}

func TestCreateURL(t *testing.T) {
	for _, test := range urlTests {
		if out := createURL(amazonBaseURL, test.recordName); out != test.expected {
			t.Errorf("output %s not equal to expected %q", out, test.expected)
		}
	}
}

func TestGetAmazonPageInfo(t *testing.T) {
	// WKM amazon url
	u := "https://www.amazon.co.uk/What-Kinda-Music-VINYL-Misch/dp/B084P38346/ref=sr_1_1?keywords=what+kinda+music+vinyl&qid=1641158805&sr=8-1"

	gotPageInfo := getAmazonPageInfo(u)
	expectedPageInfo := Record{
		Artist:      "Tom Misch & Yussef Dayes",
		Album:       "What Kinda Music",
		amazonUrl:   u,
		AmazonPrice: "£xx.xx",
	}

	if gotPageInfo.Album != expectedPageInfo.Album {
		t.Errorf("output %s not equal to expected %s", gotPageInfo.Album, expectedPageInfo.Album)
	}
	if gotPageInfo.Artist != expectedPageInfo.Artist {
		t.Errorf("output %s not equal to expected %s", gotPageInfo.Artist, expectedPageInfo.Artist)
	}
	// remove numbers to account for varying string
	re := regexp.MustCompile(`\d`)
	gotPrice := string(re.ReplaceAll([]byte(gotPageInfo.AmazonPrice), []byte("x")))

	if gotPrice != expectedPageInfo.AmazonPrice {
		t.Errorf("output %s not equal to expected %s", gotPrice, expectedPageInfo.AmazonPrice)
	}
}

func TestGetRecords(t *testing.T) {
	var sing, parr Records
	urls := readURLs("./data/input.txt")
	parr = getRecords(urls)
	for _, u := range urls {
		sing = append(
			sing,
			getAmazonPageInfo(u),
		)
	}
	if reflect.DeepEqual(sing, parr) {
		t.Error("concurrent and non-concurrent outputs do not match")
	}
}
*/
