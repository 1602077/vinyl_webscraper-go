package postgres

import (
	"reflect"
	"testing"

	r "github.com/1602077/webscraper/pkg/records"
)

func TestQueryRecordAllRows(t *testing.T) {
	db := ConnectToDB(config)
	defer db.Close()

	ExecuteFromSQLFile(db, config, "../../data/wipeTables.sql")
	ExecuteFromSQLFile(db, config, "../../data/testData.sql")

	result := QueryRecordAllRows(db)
	defer result.Close()
	if result == nil {
		t.Errorf("query failed, returned nil.")
	}
}

func TestReadQueryToRecord(t *testing.T) {
	db := ConnectToDB(config)
	defer db.Close()

	ExecuteFromSQLFile(db, config, "../../data/wipeTables.sql")
	ExecuteFromSQLFile(db, config, "../../data/testData.sql")

	result := QueryRecordAllRows(db)
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

func TestInsertIntoRecordTable(t *testing.T) {
	db := ConnectToDB(config)
	defer db.Close()

	ExecuteFromSQLFile(db, config, "../../data/wipeTables.sql")

	TestRecords := r.Records{
		r.NewRecord("Tom Misch", "Geography", "", ""),
		r.NewRecord("John Mayer", "Battle Studies", "", ""),
	}

	InsertIntoRecordTable(db, TestRecords)

	result := QueryRecordAllRows(db)
	defer result.Close()

	QueryRecords := ReadQueryToRecord(result)

	if !reflect.DeepEqual(TestRecords, QueryRecords) {
		t.Errorf("error: insert failed: Inserted records %v and returned records %v do not match.", TestRecords, QueryRecords)
	}
}
