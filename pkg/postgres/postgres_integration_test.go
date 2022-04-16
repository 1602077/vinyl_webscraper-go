package postgres

import (
	"reflect"
	"testing"
	"time"

	r "github.com/1602077/webscraper/pkg/records"
)

// Tests integration between records and postgres pkgs to confirm that records
// be inserted, read and return from db.
func TestInsertRecord(t *testing.T) {
	insertRec := r.Records{
		r.NewRecord("Tom Misch", "What Kinda Music", "", float32(25)),
		r.NewRecord("Bon Iver", "Bon Iver", "", float32(20)),
		r.NewRecord("Diana Ross", "Diana", "", float32(10)),
	}

	pg := NewPostgresCli(ENV_FILEPATH).Connect().wipe()
	for _, rec := range insertRec {
		pg.InsertRecord(rec)
	}

	results := pg.GetCurrentRecordPrices()
	returnedRec := ReadQueryToRecords(results)

	if !reflect.DeepEqual(insertRec, returnedRec) {
		t.Errorf("Records inserted do not match that returned by read operation")
	}
	pg.wipe().Close()
}

// Creates a price history in prices table and then runs GetRecordPrices to
// confirm inserted matches returned.
func TestGetRecordPrices(t *testing.T) {
	var p1, p2 float32 = 10.00, 11.50

	day := time.Date(2022, 04, 16, 0, 0, 0, 0, time.Local)
	yesterday := day.AddDate(0, 0, -1)

	rec := r.NewRecord("Chaka Khan", "I feel for you", "", p1)

	pg := NewPostgresCli(ENV_FILEPATH).Connect().wipe()
	rID, _ := pg.InsertRecord(rec) // Insert into `records` table

	insert := `
		INSERT INTO
			prices (date, price, record_id)
		VALUES
			($1, $2, $3);`
	pg.db.QueryRow(insert, day, p1, rID)
	pg.db.QueryRow(insert, yesterday, p2, rID)

	returned := pg.GetAllRecordPrices(rec)
	expected := map[string]float32{
		day.Format("2006-11-02"):       p1,
		yesterday.Format("2006-11-02"): p2,
	}

	if !reflect.DeepEqual(returned, expected) {
		t.Errorf("Pricing history returned does not matched inserted.")
	}
}
