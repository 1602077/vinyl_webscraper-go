package postgres

import (
	"log"
	"reflect"
	"testing"

	r "github.com/1602077/webscraper/pkg/records"
)

func TestQueryRecordAllRows(t *testing.T) {
	pg := NewPostgresCli(DBNAME).
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
	pg := NewPostgresCli(DBNAME).
		Connect().
		wipe().
		insertTestData()

	result := pg.QueryRecordAllRows()
	defer result.Close()

	Records := ReadQueryToRecord(result)

	length := 0
	for range Records {
		length++
	}

	if length != 3 {
		t.Errorf("expected 3 rows to be returned, got %v", length)
	}
}

func TestInsertRecord(t *testing.T) {
	pg := NewPostgresCli(DBNAME).
		Connect().
		wipe()

	TestRecords := r.Records{
		r.NewRecord("Tom Misch", "Geography", "", ""),
		r.NewRecord("John Mayer", "Battle Studies", "", ""),
	}
	for _, rr := range TestRecords {
		id := pg.InsertRecordMaster(rr)
		log.Printf("%s inserted at id %v", rr.GetAlbum(), id)
	}

	result := pg.QueryRecordAllRows()
	defer result.Close()

	QueryRecords := ReadQueryToRecord(result)

	if !reflect.DeepEqual(TestRecords, QueryRecords) {
		t.Errorf("error: insert failed: Inserted records %v and returned records %v do not match.", TestRecords, QueryRecords)
	}
}

var recThatExists = r.NewRecord("TOM MISCH", "WHAT KINDA MUSIC", "", "20")
var recThatExists2 = r.NewRecord("TOM MISCH", "WHAT KINDA MUSIC", "", "25")
var recThatDoesNotExist = r.NewRecord("Bon Iver", "i,i", "", "10")

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
	pg := NewPostgresCli(DBNAME).
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

func TestInsertRecordPricing(t *testing.T) {
	pg := NewPostgresCli(DBNAME).
		Connect().
		wipe().
		insertTestData()

	NumRecordRows := len(ReadQueryToRecord(pg.QueryRecordAllRows()))
	log.Println("Initial number of records")

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !tt.exists {
				NumRecordRows++
			}
			pg.InserRecordAllTables(tt.record)
			CurrRows := len(ReadQueryToRecord(pg.QueryRecordAllRows()))
			// if record exists should not change number of rows
			// else should increase number of rows by amount of records
			if NumRecordRows != CurrRows {
				t.Errorf("err: expected %v rows, got %v rows", NumRecordRows, CurrRows)
			}
		})
	}

	t.Run("multiple inserts only write once to pricing table if day is the same",
		func(t *testing.T) {
			pg.wipe()
			pg.InserRecordAllTables(recThatExists)
			pg.InserRecordAllTables(recThatExists2)
			pg.InserRecordAllTables(recThatDoesNotExist)

			//get length of twos; two duplicate record writes so expect pricing table to have 2 rows
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
