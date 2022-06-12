package postgres

import (
	"bytes"
	"database/sql"
	"fmt"
	"log"
	"os/exec"
	"testing"

	"github.com/1602077/webscraper/go/pkg/records"
	_ "github.com/1602077/webscraper/go/testing"
)

var (
	TEST_ENV_FILEPATH    string = "../.env.testing"
	EXAMPLE_ENV_FILEPATH string = "../.env.example"
)

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

func TestGetAllRecordsQuery(t *testing.T) {
	setup()
	defer teardown()

	result := pg.GetAllRecords()
	defer result.Close()
	if result == nil {
		t.Errorf("query failed, returned nil.")
	}
}

func TestReadRecordsTableQueryToRecord(t *testing.T) {
	setup()
	defer teardown()

	result := pg.GetAllRecords()

	Records := ReadRecordsTableQueryToRecord(result)

	if len(Records) != 3 {
		t.Errorf("expected 3 rows to be returned, got %v", len(Records))
	}
}

var recThatExists = records.NewRecord("TOM MISCH", "WHAT KINDA MUSIC", "", 20)
var recThatExists2 = records.NewRecord("TOM MISCH", "WHAT KINDA MUSIC", "", 25)
var recThatDoesNotExist = records.NewRecord("Bon Iver", "i,i", "", 10)

var tests = []struct {
	name   string
	record *records.Record
	id     int
	exists bool
}{
	{"RecordExists", recThatExists, 1, true},
	{"RecordThatDoesNotExist", recThatDoesNotExist, 0, false},
}

func TestGetRecordID(t *testing.T) {
	setupNoData()
	defer teardown()

	pg.InsertRecord(recThatExists)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			id, ok := pg.GetRecordID(tt.record)
			if tt.exists != ok {
				t.Errorf("expected %t, got %t", tt.exists, ok)
			}
			if tt.id != id {
				t.Errorf("recordIDs do not match, got %v, want %v.", id, tt.id)
			}
		})
	}
}

func TestInsertRecordPricing(t *testing.T) {
	setup()
	defer teardown()

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

var pg *PgInstance

// setup initialises testing environment with data
func setup() {
	pg = GetPgInstance().
		Connect(TEST_ENV_FILEPATH).
		wipe().
		insertTestData()
}

// setupNoData initialises testing environment with empty sql tables
func setupNoData() {
	pg = GetPgInstance().
		Connect(TEST_ENV_FILEPATH).
		wipe()
}

// teardown wipes all tables in test db and closes pg connection
func teardown() {
	pg.Close()
}

// Run command from '.sql' file on database.
func (pg *PgInstance) executeFromSQLFile(envFilename, sqlFilename string) {
	host := GetEnVar(envFilename, "DB_HOST")
	user := GetEnVar(envFilename, "DB_USER")
	dbname := GetEnVar(envFilename, "DB_NAME")
	wd := GetEnVar(envFilename, "WORKDIR")

	cmd := exec.Command("psql", "-U", user, "-h", host, "-d", dbname, "-a", "-f", wd+sqlFilename)

	var out, stderr bytes.Buffer
	cmd.Stdout, cmd.Stderr = &out, &stderr

	err := cmd.Run()
	if err != nil {
		log.Fatalf("Error executing query. Command Output: %+v\n: %+v, %v", out.String(), stderr.String(), err)
	}
}

// Clears all data from tables in pg db
func (pg *PgInstance) wipe() *PgInstance {
	pg.executeFromSQLFile(TEST_ENV_FILEPATH, "/sql/wipeTables.sql")
	return pg
}

// Insert sample data for testing
func (pg *PgInstance) insertTestData() *PgInstance {
	pg.executeFromSQLFile(TEST_ENV_FILEPATH, "/sql/testData.sql")
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
func ReadRecordsTableQueryToRecord(rows *sql.Rows) records.Records {
	var Records records.Records
	for rows.Next() {
		var id, art, alb string
		if err := rows.Scan(&id, &art, &alb); err != nil {
			break
		}
		Records = append(Records, records.NewRecord(art, alb, "", 0))
	}
	if err := rows.Err(); err != nil {
		fmt.Printf("error: query row read failed")
	}
	return Records
}
