package main

import (
	"reflect"
	"regexp"
	"testing"
)

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
