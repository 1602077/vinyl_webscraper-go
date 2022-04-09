package postgres

import (
	"reflect"
	"testing"

	r "github.com/1602077/webscraper/pkg/records"
)

// Tests integration between records and postgres pkgs to confirm that records
// be inserted, read and return from db.
func TestRecordInsertIntoDatabase(t *testing.T) {
	insertRec := r.Records{
		r.NewRecord("Tom Misch", "What Kinda Music", "", float32(25)),
		r.NewRecord("Bon Iver", "Bon Iver", "", float32(20)),
		r.NewRecord("Diana Ross", "Diana", "", float32(10)),
	}

	pg := NewPostgresCli(ENV_FILEPATH).Connect().wipe()
	for _, rec := range insertRec {
		pg.InsertRecordAllTables(rec)
	}

	results := pg.GetCurrentRecordPrices()
	returnedRec := ReadQueryToRecords(results)

	if !reflect.DeepEqual(insertRec, returnedRec) {
		t.Errorf("Records inserted do not match that returned by read operation")
	}

	pg.wipe().Close()
}
