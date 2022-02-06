package main

import (
	"fmt"
	"reflect"
	"regexp"
	"testing"
)

func TestArtistParse(t *testing.T) {
	tests := []struct{ str, name string }{
		{"Aphex Twin 678 ratings  Learn more about free returns. ...", "Aphex Twin"},
		{"Aphex Twin 67 ratings  Learn more about free returns. ...", "Aphex Twin"},
		{"Tom Misch 12 ratings ...", "Tom Misch"},
		{"0 Yussef Dayes 12345 ratings ...", "0 Yussef Dayes"},
	}

	for _, tt := range tests {
		testname := fmt.Sprintf("testing %s parse", tt.name)
		t.Run(testname, func(t *testing.T) {
			got := parseArtist(tt.str)
			want := tt.name
			if got != want {
				t.Errorf("artist parse failed: want %v, got %v", want, got)
			}
		})
	}
}

func TestPriceParse(t *testing.T) {
	tests := []struct{ str, price string }{
		{"£21.72£23.03", "£21.72"},
		{"£21.72", "£21.72"},
		{"£121.72", "£121.72"},
		{"£121.72teststrgkjg", "£121.72"},
	}

	for _, tt := range tests {
		testname := fmt.Sprintf("testing %s parse", tt.price)
		t.Run(testname, func(t *testing.T) {
			got := parsePrice(tt.str)
			want := tt.price
			if got != want {
				t.Errorf("price parse failed: want %v, got %v", want, got)
			}
		})
	}
}

func TestGetAmazonPageInfo(t *testing.T) {
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

	// remove numbers to account for varying price
	re := regexp.MustCompile(`\d`)
	gotPrice := string(re.ReplaceAll([]byte(gotPageInfo.AmazonPrice), []byte("x")))
	if gotPrice != "" && gotPrice != expectedPageInfo.AmazonPrice {
		t.Errorf("output %s not equal to expected %s", gotPrice, expectedPageInfo.AmazonPrice)
	}
}

// Tests that concurrent implimentation matches single threaded version
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
		t.Errorf("non-concurrent and concurrentoutputs do not match.\nexpected: %v.\ngot:%v.", sing, parr)
	}
}
