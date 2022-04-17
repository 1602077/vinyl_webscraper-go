// api methods to read & write to postgres database
package postgres

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	r "github.com/1602077/webscraper/pkg/records"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type PgInstance struct {
	db *sql.DB
}

var pginstance *PgInstance

// NewPgInstace() is a factory function for creating a singleton PgInstance
func GetPgInstance() *PgInstance {
	if pginstance == nil {
		pginstance = new(PgInstance)
	}
	return pginstance
}

// GetEnVar uses godot to read env variables from a .env file
func GetEnVar(filepath, key string) string {
	err := godotenv.Load(filepath)
	if err != nil {
		log.Fatalf("Error loading .env file: %s", err)
	}
	return os.Getenv(key)
}

// Opens a connection to database specified by .env file
func (pg *PgInstance) Connect(filepath string) *PgInstance {

	host := GetEnVar(filepath, "DB_HOST")
	port, err := strconv.Atoi(GetEnVar(filepath, "DB_PORT"))
	if err != nil {
		log.Fatalf("Connect() failed: Port convervsion failed: %v\n", err)
	}
	user := GetEnVar(filepath, "DB_USER")
	pwd := GetEnVar(filepath, "DB_PASSWORD")
	dbname := GetEnVar(filepath, "DB_NAME")

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, pwd, dbname)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatalf("err: opening connection to database '%s' failed.", dbname)
	}

	err = db.Ping()
	if err != nil {
		log.Fatalf("err: ping to database '%s' failed: %s", dbname, err)
	}

	log.Printf("connection to database '%s' successfully opened.\n", dbname)
	pg.db = db
	return pg
}

// Closes connection to database
func (pg *PgInstance) Close() {
	pg.db.Close()
	log.Print("closing to database closed.")
}

// Retrieves the id of a record from 'records' table
func (pg *PgInstance) GetRecordID(rec *r.Record) (int, bool) {
	existsQuery := `
		SELECT id
		FROM records
		WHERE artist = $1 AND album = $2
		LIMIT 1;`

	var recordID int
	if err := pg.db.QueryRow(existsQuery, rec.GetArtist(), rec.GetAlbum()).Scan(&recordID); err == sql.ErrNoRows {
		return 0, false
	}
	return recordID, true
}

// Retrieves the id of a price row for a given record_id and date
func (pg *PgInstance) GetPriceID(recordID int, date time.Time) (int, bool) {
	existsQuery := `
		SELECT id
		FROM prices
		WHERE date = $1 AND record_id = $2
		LIMIT 1;`

	var priceID int
	if err := pg.db.QueryRow(existsQuery, date, recordID).Scan(&priceID); err == sql.ErrNoRows {
		return 0, false
	}
	return priceID, true
}

// Gets most recent prices of all records
func (pg *PgInstance) GetCurrentRecordPrices() *sql.Rows {
	rows, err := pg.db.Query(`
		SELECT r.Artist, r.Album, p.MaxPrice
		FROM records r
		INNER JOIN (
			SELECT record_id, MAX(Date) as MaxDate, MAX(price) as MaxPrice
			FROM prices
			GROUP BY record_id
		) p ON p.record_id = r.id;`)

	if err != nil {
		log.Fatalf("err: GetCurrentRecordPrices() failed: %v.", err)
	}
	return rows
}

func (pg *PgInstance) GetAllRecordPrices(r *r.Record) map[string]float32 {
	rows, err := pg.db.Query(`
		SELECT date, price
		FROM prices
		WHERE record_id IN (
			SELECT id
			FROM records
			WHERE album = $1 AND artist = $2
		);`, r.GetAlbum(), r.GetArtist())

	if err != nil {
		log.Fatalf("err: GetRecordPrices(%s) failed: %v.", r.GetAlbum(), err)
	}

	prices := make(map[string]float32)
	for rows.Next() {
		var date time.Time
		var price float32
		if err := rows.Scan(&date, &price); err != nil {
			break
		}
		datestring := date.Format("2006-11-02")
		prices[datestring] = price
	}
	return prices
}

// Checks if record exists in 'records' table and adds if not & then inserts into pricing table
func (pg *PgInstance) InsertRecord(rec *r.Record) (int, int) {
	recordID, ok := pg.GetRecordID(rec)
	if !ok {
		insertQuery := `
			INSERT INTO
				records (artist, album)
			VALUES
				($1, $2)
			RETURNING ID;`

		var rID int
		if err := pg.db.QueryRow(insertQuery, rec.GetArtist(), rec.GetAlbum()).Scan(&rID); err != nil {
			log.Fatalf("err: InsertRecordIntoRecords() failed: %v.", err)
		}
		recordID = rID
	}

	today := time.Now()
	priceID, ok := pg.GetPriceID(recordID, today)
	if ok {
		updateQuery := `
			UPDATE prices
			SET price = $1
			WHERE date = $2 AND record_id = $3
			RETURNING ID;`

		if err := pg.db.QueryRow(updateQuery, rec.GetPrice(), today, recordID).Scan(&priceID); err == sql.ErrNoRows {
			return recordID, priceID
		}
		log.Printf("%s: updated in db.", rec.GetAlbum())
		return recordID, priceID
	}

	insertQuery := `
		INSERT INTO
			prices (date, price, record_id)
		VALUES
			($1, $2, $3)
		RETURNING ID;`

	pg.db.QueryRow(insertQuery, today, rec.GetPrice(), recordID).Scan(&priceID)
	log.Printf("%s: written to db.", rec.GetAlbum())
	return recordID, priceID
}

func ReadQueryToRecords(rows *sql.Rows) r.Records {
	var Records r.Records
	for rows.Next() {
		var art, alb string
		var price float32
		if err := rows.Scan(&art, &alb, &price); err != nil {
			break
		}
		Records = append(Records, r.NewRecord(art, alb, "", price))
	}
	if err := rows.Err(); err != nil {
		log.Fatalf("err: ReadQueryToRecords() failed: %v.", err)
	}
	return Records
}

func (pg *PgInstance) PrintCurrentRecordPrices() {
	rows := pg.GetCurrentRecordPrices()
	rec := ReadQueryToRecords(rows)
	rec.PrintRecords()
}
