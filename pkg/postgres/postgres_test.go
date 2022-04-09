package postgres

import (
	"log"
	"reflect"
	"testing"

	r "github.com/1602077/webscraper/pkg/records"
)

var ENV_FILEPATH string = "../../.env.testing"

func TestQueryRecordAllRows(t *testing.T) {
	pg := NewPostgresCli(ENV_FILEPATH).
		Connect().
		wipe().
		insertTestData()

	result := pg.QueryRecordAllRows()
	defer result.Close()
	if result == nil {
		t.Errorf("query failed, returned nil.")
	}
}

func TestReadQueryToRecord(t *testing.T) {
	pg := NewPostgresCli(ENV_FILEPATH).
		Connect().
		wipe().
		insertTestData()

	result := pg.QueryRecordAllRows()
	defer result.Close()

	Records := ReadPartialQueryToRecord(result)

	length := 0
	for range Records {
		length++
	}

	if length != 3 {
		t.Errorf("expected 3 rows to be returned, got %v", length)
	}
}

func TestInsertRecord(t *testing.T) {
	pg := NewPostgresCli(ENV_FILEPATH).
		Connect().
		wipe()

	TestRecords := r.Records{
		r.NewRecord("Tom Misch", "Geography", "", 0),
		r.NewRecord("John Mayer", "Battle Studies", "", 0),
	}
	for _, rr := range TestRecords {
		id := pg.InsertRecordMaster(rr)
		log.Printf("%s inserted at id %v", rr.GetAlbum(), id)
	}

	result := pg.QueryRecordAllRows()
	defer result.Close()

	QueryRecords := ReadPartialQueryToRecord(result)

	if !reflect.DeepEqual(TestRecords, QueryRecords) {
		t.Errorf("error: insert failed: Inserted records %v and returned records %v do not match.", TestRecords, QueryRecords)
	}
}

var recThatExists = r.NewRecord("TOM MISCH", "WHAT KINDA MUSIC", "", 20)
var recThatExists2 = r.NewRecord("TOM MISCH", "WHAT KINDA MUSIC", "", 25)
var recThatDoesNotExist = r.NewRecord("Bon Iver", "i,i", "", 10)

var tests = []struct {
	name   string
	record *r.Record
	id     int
	exists bool
}{
	{"RecordExists", recThatExists, 1, true},
	{"RecordThatDoesNotExist", recThatDoesNotExist, 0, false},
}

func TestGetRecordID(t *testing.T) {
	pg := NewPostgresCli(ENV_FILEPATH).
		Connect().
		wipe()

	pg.InsertRecordMaster(recThatExists)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			id, ok := pg.GetRecordID(tt.record)
			if tt.exists != ok {
				t.Errorf("err: expected %t, got %t", tt.exists, ok)
			}
			if tt.id != id {
				t.Errorf("err: recordIDs do not match, got %v, want %v.", id, tt.id)
			}
		})
	}
}

/*
func TestInsertRecordPricing(t *testing.T) {
	pg := NewPostgresCli(ENV_FILEPATH).
		Connect().
		wipe().
		insertTestData()

	NumRecordRows := len(ReadPartialQueryToRecord(pg.QueryRecordAllRows()))

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !tt.exists {
				NumRecordRows++
			}
			pg.InsertRecordAllTables(tt.record)
			CurrRows := len(ReadPartialQueryToRecord(pg.QueryRecordAllRows()))
			// if record exists should not change number of rows
			// else should increase number of rows by amount of records
			if NumRecordRows != CurrRows {
				t.Errorf("err: expected %v rows, got %v rows", NumRecordRows, CurrRows)
			}
		})
	}

	t.Run("multiple inserts into db only write once to pricing table if day is the same",
		func(t *testing.T) {
			pg.wipe()
			pg.InsertRecordAllTables(recThatExists)
			pg.InsertRecordAllTables(recThatExists2)
			pg.InsertRecordAllTables(recThatDoesNotExist)

			// two duplicate record writes so expect pricing table to have only 2 rows
			var numRows int
			rows := pg.QueryPriceAllRows()
			for rows.Next() {
				numRows++
			}
			if numRows != 2 {
				t.Errorf("err: expected 2 rows, got %v", numRows)
			}
		})
}
*/
