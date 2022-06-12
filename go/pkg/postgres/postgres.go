// api methods to read & write to postgres database
package postgres

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/1602077/webscraper/go/pkg/records"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type PgInstance struct {
	db *sql.DB
}

var pginstance *PgInstance

// NewPgInstace is a factory function for creating a singleton PgInstance
// TODO:  Embed this into the Connect method and remove redundant PgInstance struct
func GetPgInstance() *PgInstance {
	if pginstance == nil {
		pginstance = new(PgInstance)
	}
	return pginstance
}

// GetEnVar uses godotenv to read in env variables specified by key from a .env filepath.
func GetEnVar(filepath, key string) string {
	err := godotenv.Load(filepath)
	if err != nil {
		log.Fatalf("error loading .env file: %s", err)
	}
	return os.Getenv(key)
}

// Connect to database specified by .env file.
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

	if err = db.Ping(); err != nil {
		log.Fatalf("err: ping to database '%s' failed: %s", dbname, err)
	}

	log.Printf("connection to database '%s' opened.\n", dbname)
	pg.db = db
	return pg
}

// Close connection to postgres database.
func (pg *PgInstance) Close() {
	pg.db.Close()
	log.Print("connection to database closed.")
}

// GetRecordID retrieves the id of the input record from 'records' table.
func (pg *PgInstance) GetRecordID(rec *records.Record) (int, bool) {
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

// GetPriceID retrieves the id of a price row for a given record_id and date.
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

// GetCurrentRecordPrice gets most recent prices of all records in pg database.
func (pg *PgInstance) GetCurrentRecordPrices() records.Records {
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

	var Records records.Records
	for rows.Next() {
		var art, alb string
		var price float32
		if err := rows.Scan(&art, &alb, &price); err != nil {
			break
		}
		Records = append(Records, records.NewRecord(art, alb, "", price))
	}
	if err := rows.Err(); err != nil {
		log.Fatalf("err: GetCurrentRecordPrices() failed: %v.", err)
	}
	return Records
}

// GetAllRecordPrices retrieves the full price history of a single input record.
func (pg *PgInstance) GetAllRecordPrices(r *records.Record) map[string]float32 {
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

// InsertRecord adds record to the 'records' table if it does not exist and
// inserts into current price into the pricing table. If a price already exists
// for the date of insert it is updated instead.
func (pg *PgInstance) InsertRecord(rec *records.Record) (int, int) {
	recordID, ok := pg.GetRecordID(rec)
	if !ok {
		insertQuery := `
			INSERT INTO
				records (artist, album)
			VALUES
				($1, $2)
			RETURNING ID;`

		var rID int
		err := pg.db.QueryRow(insertQuery, rec.GetArtist(), rec.GetAlbum()).Scan(&rID)
		if err != nil {
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

		err := pg.db.QueryRow(updateQuery, rec.GetPrice(), today, recordID).Scan(&priceID)
		if err == sql.ErrNoRows {
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

// PrintCurrentPrices prints the artist, album and most recent price for
// all records in database as tab written table.
func (pg *PgInstance) PrintCurrentPrices() {
	rec := pg.GetCurrentRecordPrices()
	rec.Print()
}

// GetRecordPriceHistory retrieves the artist, album and full price history for
// the record specified by the input id.
func (pg *PgInstance) GetRecordPriceHistory(id int) *records.RecordPriceHistory {
	rIdQuery := `
		SELECT r.artist, r.album
		FROM records r
		WHERE r.id = $1;`

	var artist, album string
	if err := pg.db.QueryRow(rIdQuery, id).Scan(&artist, &album); err != nil {
		if err == sql.ErrNoRows {
			log.Printf("err: GetRecordPriceHistory: no record id found at id %v\n", id)
			return nil
		}
		log.Printf("GetRecordPriceHistory: parsing sql rId query failed: %s\n", err)
	}

	phQuery := `
		SELECT p.date, p.price
		FROM prices p
		WHERE p.record_id = $1
		ORDER BY p.date ASC;`

	rows, err := pg.db.Query(phQuery, id)
	if err != nil {
		log.Printf("err: GetRecordPriceHistory: price history query failed: %s\n", err)
	}

	var priceHistory []*records.PriceHist
	for rows.Next() {
		var date string
		var price float32
		if err := rows.Scan(&date, &price); err != nil {
			break
		}
		priceHistory = append(priceHistory, &records.PriceHist{Date: date, Price: price})
	}

	return &records.RecordPriceHistory{
		Id:           id,
		Artist:       artist,
		Album:        album,
		PriceHistory: priceHistory,
	}
}
