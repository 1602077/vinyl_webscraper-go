package postgres

import (
	"bytes"
	"database/sql"
	"fmt"
	"log"
	"os"
	"os/exec"
	"time"

	r "github.com/1602077/webscraper/pkg/records"
	_ "github.com/lib/pq"
)

type pgConfig struct {
	host, user, password, dbname string
	port                         int
}

var config = pgConfig{
	host:     "localhost",
	port:     5432,
	user:     os.Getenv("DBUSER"),
	password: os.Getenv("DBPASS"),
	dbname:   "records",
}

/*
type database struct {
	db *sql.DB
	pgConfig
}

var db = database{
	db:       ConnectToDB(config),
	pgConfig: config,
}
*/

func ConnectToDB(dbConfig pgConfig) *sql.DB {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		dbConfig.host, dbConfig.port, dbConfig.user, dbConfig.password, dbConfig.dbname)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatalf("err: opening connection to database '%s' failed.", dbConfig.dbname)
	}

	err = db.Ping()
	if err != nil {
		log.Fatalf("err: ping to database '%s' failed.", dbConfig.dbname)
	}

	log.Printf("connection to database '%s' successfully opened.", dbConfig.dbname)

	return db
}

func ExecuteFromSQLFile(db *sql.DB, c pgConfig, filename string) {
	cmd := exec.Command("psql", "-U", c.user, "-h", c.host, "-d", c.dbname, "-a", "-f", filename)

	var out, stderr bytes.Buffer
	cmd.Stdout, cmd.Stderr = &out, &stderr

	err := cmd.Run()
	if err != nil {
		log.Fatalf("Error executing query. Command Output: %+v\n: %+v, %v", out.String(), stderr.String(), err)
	}
}

// Runs "SELECT * FROM record"
func QueryRecordAllRows(db *sql.DB) *sql.Rows {
	rows, err := db.Query("SELECT * FROM record")
	if err != nil {
		log.Fatal("err: quering all rows on 'records' table failed.")
	}
	return rows
}

// Reads in the result of a db.Query(...) [*sql.Rows] to r.Records type
func ReadQueryToRecord(rows *sql.Rows) r.Records {
	var Records r.Records
	for rows.Next() {
		var id, art, alb string
		if err := rows.Scan(&id, &art, &alb); err != nil {
			break
		}
		Records = append(Records, r.NewRecord(art, alb, "", ""))
	}
	if err := rows.Err(); err != nil {
		fmt.Printf("error: query row read failed")
	}
	return Records
}

// Inserts a single record into 'record' table and returns it's id
func InsertRecordMaster(db *sql.DB, rec *r.Record) int {
	insertQuery := `
		INSERT INTO
			record (artist, album)
		VALUES
			($1, $2)
		RETURNING ID;`

	var id int
	if err := db.QueryRow(insertQuery, rec.GetArtist(), rec.GetAlbum()).Scan(&id); err != nil {
		log.Print("err: insert into 'record' table failed.")
	}
	return id
}

// Retrieves the id of a record from 'record' table
func GetRecordID(db *sql.DB, rec *r.Record) (int, bool) {
	existsQuery := `
		SELECT id
		FROM record
		WHERE artist = $1 AND album = $2
		LIMIT 1;`

	var id int
	if err := db.QueryRow(existsQuery, rec.GetArtist(), rec.GetAlbum()).Scan(&id); err != nil {
		log.Printf("err: record not found in 'record' table.")
		return 0, false
	}
	return id, true
}

// Retrieves the id of a price row for a given record_id and date
func GetPriceID(db *sql.DB, recordID int, date time.Time) (int, bool) {
	existsQuery := `
		SELECT id
		FROM price
		WHERE date = $1 AND record_id = $2
		LIMIT 1;`

	var id int
	if err := db.QueryRow(existsQuery, date, recordID).Scan(&id); err != nil {
		log.Printf("err: price not found in 'price' table.")
		return 0, false
	}
	return id, true
}

// Checks if record exists in 'record' table and adds if not & then inserts into pricing table
func InsertRecordPricing(db *sql.DB, rec *r.Record) int {
	// check if exists in 'record table'
	recordID, ok := GetRecordID(db, rec)
	if !ok {
		recordID = InsertRecordMaster(db, rec)
	}

	today := time.Now()
	// today := time.Now().Format("2006-01-02")

	priceID, ok := GetPriceID(db, recordID, today)
	if ok {
		// replace instead
		updateQuery := `
			UPDATE price
			SET price = $1
			WHERE date = $2 AND record_id = $3
			RETURNING ID;`

		if err := db.QueryRow(updateQuery, rec.GetPrice(), today, recordID).Scan(&priceID); err != nil {
			log.Printf("err:  price not found in 'price' table.")
			return priceID
		}
	}

	insertQuery := `
		INSERT INTO
			price (date, price, record_id)
		VALUES
			($1, $2, $3)
		RETURNING ID;`

	if err := db.QueryRow(insertQuery, today, rec.GetPrice(), recordID).Scan(&priceID); err != nil {
		log.Printf("err: insert into 'price' table failed.")
	}
	return priceID
}

// Runs "SELECT * FROM price"
func QueryPriceAllRows(db *sql.DB) *sql.Rows {
	rows, err := db.Query("SELECT * FROM price")
	if err != nil {
		log.Fatal("err: quering all rows on 'price' table failed.")
	}
	return rows
}
