package postgres

import (
	"bytes"
	"database/sql"
	"fmt"
	"log"
	"os/exec"
	"testing"

	r "github.com/1602077/webscraper/pkg/records"
	_ "github.com/1602077/webscraper/testing"
)

var TEST_ENV_FILEPATH string = ".env.testing"
var EXAMPLE_ENV_FILEPATH string = ".env.example"

func TestGetPgInstace(t *testing.T) {
	pginstance1 := GetPgInstance()
	if pginstance1 == nil {
		t.Error("Expected pointer to singleton after calling GetPgInstance(), got nil")
	}

	pginstance2 := GetPgInstance()
	if pginstance1 != pginstance2 {
		t.Errorf("Expected same instance in pginstance2, but got a different one")
	}
}

func TestGetEnVar(t *testing.T) {
	key, value := "EXAMPLE_KEY", "EXAMPLE_VALUE"
	actual := GetEnVar(EXAMPLE_ENV_FILEPATH, key)
	if actual != value {
		t.Errorf("GetEnVar(%s): expected %s, actual %s", key, value, actual)
	}
}

func TestQueryRecordAllRows(t *testing.T) {
	pg := GetPgInstance().
		Connect(TEST_ENV_FILEPATH).
		wipe().
		insertTestData()

	result := pg.GetAllRecords()
	defer result.Close()
	if result == nil {
		t.Errorf("query failed, returned nil.")
	}
}

func TestGetAllRecords(t *testing.T) {
	pg := GetPgInstance().
		Connect(TEST_ENV_FILEPATH).
		wipe().
		insertTestData()

	result := pg.GetAllRecords()

	Records := ReadRecordsTableQueryToRecord(result)

	if len(Records) != 3 {
		t.Errorf("expected 3 rows to be returned, got %v", len(Records))
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
	pg := GetPgInstance().
		Connect(TEST_ENV_FILEPATH).
		wipe()

	pg.InsertRecord(recThatExists)

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
	pg := GetPgInstance().
		Connect(TEST_ENV_FILEPATH).
		wipe().
		insertTestData()

	NumRecordRows := len(ReadRecordsTableQueryToRecord(pg.GetAllRecords()))

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !tt.exists {
				NumRecordRows++
			}
			pg.InsertRecord(tt.record)
			CurrRows := len(ReadRecordsTableQueryToRecord(pg.GetAllRecords()))
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
			rows := pg.GetAllPrices()
			for rows.Next() {
				numRows++
			}
			if numRows != 2 {
				t.Errorf("err: expected 2 rows, got %v", numRows)
			}
		})
}

// Run command from '.sql' file on database.
func (pg *PgInstance) executeFromSQLFile(envFilename, sqlFilename string) {
	host := GetEnVar(envFilename, "DB_HOST")
	user := GetEnVar(envFilename, "DB_USER")
	dbname := GetEnVar(envFilename, "DB_NAME")

	cmd := exec.Command("psql", "-U", user, "-h", host, "-d", dbname, "-a", "-f", sqlFilename)

	var out, stderr bytes.Buffer
	cmd.Stdout, cmd.Stderr = &out, &stderr

	err := cmd.Run()
	if err != nil {
		log.Fatalf("Error executing query. Command Output: %+v\n: %+v, %v", out.String(), stderr.String(), err)
	}
}

// Clears all data from tables in pg db
func (pg *PgInstance) wipe() *PgInstance {
	pg.executeFromSQLFile(TEST_ENV_FILEPATH, "./data/wipeTables.sql")
	return pg
}

// Insert sample data for testing
func (pg *PgInstance) insertTestData() *PgInstance {
	pg.executeFromSQLFile(TEST_ENV_FILEPATH, "./data/testData.sql")
	return pg
}

// Runs "SELECT * FROM records"
func (pg *PgInstance) GetAllRecords() *sql.Rows {
	rows, err := pg.db.Query("SELECT * FROM records;")
	if err != nil {
		log.Fatalf("err: QueryRecordAllRows() failed: %v.", err)
	}
	return rows
}

// Runs "SELECT * FROM prices"
func (pg *PgInstance) GetAllPrices() *sql.Rows {
	rows, err := pg.db.Query("SELECT * FROM prices;")
	if err != nil {
		log.Fatalf("err: QueryPriceAllRows() failed: %v", err)
	}
	return rows
}

// Reads in the result of a db.Query(...) [*sql.Rows] to r.Records type
func ReadRecordsTableQueryToRecord(rows *sql.Rows) r.Records {
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