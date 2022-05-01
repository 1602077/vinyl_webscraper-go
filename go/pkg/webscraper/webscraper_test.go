package webscraper

import (
	"fmt"
	"reflect"
	"strconv"
	"testing"

	db "github.com/1602077/webscraper/go/pkg/postgres"
	r "github.com/1602077/webscraper/go/pkg/records"
	_ "github.com/1602077/webscraper/go/testing"
)

var ENV_FILEPATH string = "../.env.testing"

func TestArtistParse(t *testing.T) {
	tests := []struct{ str, name string }{
		{"Aphex Twin 678 ratings  Learn more about free returns. ...", "Aphex Twin"},
		{"Aphex Twin 67 ratings  Learn more about free returns. ...", "Aphex Twin"},
		{"Tom Misch 12 ratings ...", "Tom Misch"},
		{"0 Yussef Dayes 12345 ratings ...", "0 Yussef Dayes"},
		{"Arctic Monkeys 6,866 ratings  Learn more about free returns.", "Arctic Monkeys"},
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
		{"£21.72£23.03", "21.72"},
		{"£21.72", "21.72"},
		{"£121.72", "121.72"},
		{"£121.72teststrgkjg", "121.72"},
	}

	for _, tt := range tests {
		testname := fmt.Sprintf("testing %s parse", tt.price)
		t.Run(testname, func(t *testing.T) {
			got := parsePrice(tt.str)
			pf, _ := strconv.ParseFloat(tt.price, 32)
			want := float32(pf)
			if got != want {
				t.Errorf("price parse failed: want %v, got %v", want, got)
			}
		})
	}
}

func TestGetAmazonPageInfo(t *testing.T) {
	u := "https://www.amazon.co.uk/AM-VINYL-Arctic-Monkeys/dp/B00DKY4NBA/ref=sr_1_4?crid=EIQTUGWC5AAR&keywords=vinyl&qid=1645263030&sprefix=vinyl%2Caps%2C83&sr=8-4"

	gotPageInfo := getAmazonPageInfo(u)
	expectedPageInfo := r.NewRecord("Arctic Monkeys", "AM", u, 0.0)
	fmt.Print(gotPageInfo)

	if gotPageInfo.GetAlbum() != expectedPageInfo.GetAlbum() {
		t.Errorf("output %s not equal to expected %s", gotPageInfo.GetAlbum(), expectedPageInfo.GetAlbum())
	}

	if gotPageInfo.GetArtist() != expectedPageInfo.GetArtist() {
		t.Errorf("output %s not equal to expected %s", gotPageInfo.GetArtist(), expectedPageInfo.GetArtist())
	}
}

// Tests that concurrent implimentation matches single threaded version
func TestGetRecords(t *testing.T) {
	wd := db.GetEnVar(ENV_FILEPATH, "WORKDIR")
	urls := ReadURLs(wd + "/input.txt")

	var sing, parr r.Records
	parr = GetRecords(urls)
	for _, u := range urls {
		sing = append(
			sing,
			getAmazonPageInfo(u),
		)
	}
	if reflect.DeepEqual(sing, parr) {
		t.Errorf("non-concurrent and concurrent outputs do not match.\nexpected: %v.\ngot:%v.", sing, parr)
	}
}
