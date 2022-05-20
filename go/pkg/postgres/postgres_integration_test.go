package postgres

import (
	"reflect"
	"testing"
	"time"

	r "github.com/1602077/webscraper/go/pkg/records"
)

// Tests integration between records and postgres pkgs to confirm that records
// be inserted, read and return from db.
func TestInsertRecord(t *testing.T) {
	setupNoData()
	defer teardown()

	insertRec := r.Records{
		r.NewRecord("Tom Misch", "What Kinda Music", "", float32(25)),
		r.NewRecord("Bon Iver", "Bon Iver", "", float32(20)),
		r.NewRecord("Diana Ross", "Diana", "", float32(10)),
	}

	for _, rec := range insertRec {
		pg.InsertRecord(rec)
	}

	results := pg.GetCurrentRecordPrices()
	returnedRec := ReadQueryToRecords(results)

	if !reflect.DeepEqual(insertRec, returnedRec) {
		t.Errorf("Records inserted do not match that returned by read operation")
	}
}

// Creates a price history in prices table and then runs GetRecordPrices to
// confirm inserted matches returned.
func TestGetRecordPrices(t *testing.T) {
	setupNoData()
	defer teardown()

	var p1, p2 float32 = 10.00, 11.50
	day1, day2 := time.Now(), time.Date(2022, 04, 16, 0, 0, 0, 0, time.Local)
	r1 := r.NewRecord("Chaka Khan", "I feel for you", "", p1)

	// Intial insert into records and prices
	pg.InsertRecord(r1)
	// Second insert into prices
	pg.db.QueryRow(`INSERT INTO prices (date, price, record_id) VALUES ($1, $2, $3);`,
		day2, p2, 1)

	returned := pg.GetAllRecordPrices(r1)
	expected := map[string]float32{
		day1.Format("2006-11-02"): p1,
		day2.Format("2006-11-02"): p2,
	}

	if !reflect.DeepEqual(returned, expected) {
		t.Errorf("Pricing history returned does not matched inserted.")
	}
}
