package postgres

import (
	"bytes"
	"database/sql"
	"fmt"
	"log"
	"os/exec"
	"testing"

	r "github.com/1602077/webscraper/pkg/records"
)

// Run command from '.sql' file on database.
func (pg *PgInstance) executeFromSQLFile(filename string) {
	cmd := exec.Command("psql", "-U", pg.config.user, "-h", pg.config.host, "-d", pg.config.dbname, "-a", "-f", filename)

	var out, stderr bytes.Buffer
	cmd.Stdout, cmd.Stderr = &out, &stderr

	err := cmd.Run()
	if err != nil {
		log.Fatalf("Error executing query. Command Output: %+v\n: %+v, %v", out.String(), stderr.String(), err)
	}
}

// Clears all data from tables in pg db
func (pg *PgInstance) wipe() *PgInstance {
	pg.executeFromSQLFile("../../data/wipeTables.sql")
	return pg
}

// Insert sample data for testing
func (pg *PgInstance) insertTestData() *PgInstance {
	pg.executeFromSQLFile("../../data/testData.sql")
	return pg
}

// Runs "SELECT * FROM records"
func (pg *PgInstance) QueryRecordAllRows() *sql.Rows {
	rows, err := pg.db.Query("SELECT * FROM records;")
	if err != nil {
		log.Fatalf("err: QueryRecordAllRows() failed: %v.", err)
	}
	return rows
}

// Reads in the result of a db.Query(...) [*sql.Rows] to r.Records type
func ReadPartialQueryToRecord(rows *sql.Rows) r.Records {
	var Records r.Records
	for rows.Next() {
		var id, art, alb string
		if err := rows.Scan(&id, &art, &alb); err != nil {
			break
		}
		Records = append(Records, r.NewRecord(art, alb, "", 0))
	}
	if err := rows.Err(); err != nil {
		fmt.Printf("error: query row read failed")
	}
	return Records
}

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

	pg.InsertRecordIntoRecords(recThatExists)

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
			pg.InsertRecord(tt.record)
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
			pg.InsertRecord(recThatExists)
			pg.InsertRecord(recThatExists2)
			pg.InsertRecord(recThatDoesNotExist)

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
